package session

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/streamingfast/dsession"
	pbworker "github.com/streamingfast/worker-pool-protocol/pb/sf/worker/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Register registers the TGM session pool with dauth
func Register() {
	dsession.Register("tgm", func(config string, logger *zap.Logger) (dsession.SessionPool, error) {
		configExpanded := os.ExpandEnv(config)

		c, err := newConfig(configExpanded)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config string %s: %w", config, err)
		}
		return newTGMSessionPool(c, logger)
	})
}

type tgmSessionPool struct {
	config                 *Config
	logger                 *zap.Logger
	remoteWorkerPoolClient pbworker.WorkerPoolClient
	conn                   *grpc.ClientConn
	sessionClosers         map[string]chan struct{}
	sessionMutex           sync.Mutex
}

type borrowedSession struct {
	done chan struct{}
}

func newTGMSessionPool(config *Config, logger *zap.Logger) (dsession.SessionPool, error) {
	logger = logger.Named("tgm-session-pool")

	// Create gRPC connection
	var opts []grpc.DialOption
	if config.Plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := &tls.Config{}
		if config.Insecure {
			tlsConfig.InsecureSkipVerify = true
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

	// Add unary interceptor for X-Api-Key header if API key is provided
	if config.IndexerApiKey != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(createApiKeyInterceptor(config.IndexerApiKey)))
	}

	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session endpoint %s: %w", config.Endpoint, err)
	}

	client := pbworker.NewWorkerPoolClient(conn)

	pool := &tgmSessionPool{
		config:                 config,
		logger:                 logger,
		remoteWorkerPoolClient: client,
		conn:                   conn,
		sessionClosers:         make(map[string]chan struct{}),
		sessionMutex:           sync.Mutex{},
	}

	return pool, nil
}

func (t *tgmSessionPool) Get(ctx context.Context, serviceName string, userID string, apiKeyID string, traceID string, onError func(error)) (string, error) {
	resp, err := t.remoteWorkerPoolClient.BorrowWorker(ctx,
		&pbworker.BorrowWorkerRequest{
			Service:  serviceName,
			UserId:   userID,
			ApiKeyId: apiKeyID,
			TraceId:  traceID,
		},
		grpc.WaitForReady(false),
	)

	if err != nil {
		// Map gRPC errors to dsession errors
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Unavailable:
				return "", fmt.Errorf("%w: %s", dsession.ErrUnavailable, grpcErr.Message())
			case codes.PermissionDenied:
				return "", fmt.Errorf("%w: %s", dsession.ErrPermissionDenied, grpcErr.Message())
			case codes.ResourceExhausted:
				return "", fmt.Errorf("%w: %s", dsession.ErrQuotaExceeded, grpcErr.Message())
			}
		}
		return "", fmt.Errorf("failed to borrow session: %w", err)
	}

	key := resp.WorkerKey
	workerStatus := resp.Status

	if workerStatus == pbworker.BorrowWorkerResponse_resource_exhausted {
		t.logger.Info("worker pool is exhausted", zap.String("worker_key", key), zap.String("status", workerStatus.String()))
		return "", fmt.Errorf("worker pool exhausted: %w", dsession.ErrConcurrentStreamLimitExceeded)
	}

	// Start keep-alive for borrowed workers
	if workerStatus == pbworker.BorrowWorkerResponse_borrowed {

		done := make(chan struct{})
		t.sessionMutex.Lock()
		t.sessionClosers[key] = done
		t.sessionMutex.Unlock()

		startKeepAlive(ctx, t.config.RequestKeepAliveDelay, done, t.remoteWorkerPoolClient, key, apiKeyID, onError, t.logger)
	}

	t.logger.Debug("borrowed request worker", zap.String("worker_key", key))

	return key, nil
}

func (t *tgmSessionPool) Release(sessionKey string) {
	t.sessionMutex.Lock()
	done := t.sessionClosers[sessionKey]
	if done != nil {
		close(done)
		delete(t.sessionClosers, sessionKey)
	}
	t.sessionMutex.Unlock()

	resp, err := t.remoteWorkerPoolClient.ReturnWorker(context.Background(),
		&pbworker.ReturnWorkerRequest{
			WorkerKey:                 sessionKey,
			MinimalWorkerLifeDuration: durationpb.New(t.config.DefaultMinimalWorkerLifeDuration),
		},
		grpc.WaitForReady(false),
	)

	t.logger.Debug("returned request worker", zap.String("key", sessionKey), zap.Stringer("status", resp.Status), zap.Error(err))
}

// createApiKeyInterceptor creates a gRPC unary interceptor that adds the X-Api-Key header
func createApiKeyInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add the X-Api-Key header to the outgoing metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Api-Key", apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// startKeepAlive starts the keep-alive goroutine for a borrowed session
func startKeepAlive(ctx context.Context, delay time.Duration, done <-chan struct{}, client pbworker.WorkerPoolClient, workerKey, apiKeyID string, onError func(error), logger *zap.Logger) {
	originalDelay := delay

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-time.After(delay):
				delay = originalDelay
				_, err := client.KeepAlive(
					ctx,
					&pbworker.KeepAliveRequest{
						WorkerKey: workerKey,
						ApiKeyId:  apiKeyID,
					},
					grpc.WaitForReady(false),
				)
				if err != nil {
					logger.Error("failed to call keep request worker alive", zap.String("worker_id", workerKey), zap.Error(err))
					if onError != nil {
						// Map gRPC errors to dsession errors
						if grpcErr, ok := status.FromError(err); ok {
							switch grpcErr.Code() {
							case codes.PermissionDenied:
								onError(fmt.Errorf("%w: %s", dsession.ErrPermissionDenied, grpcErr.Message()))
								return
							case codes.ResourceExhausted:
								onError(fmt.Errorf("%w: %s", dsession.ErrQuotaExceeded, grpcErr.Message()))
								return
							}
						}
						delay = time.Second
						logger.Info("keep-alive failed, retrying", zap.String("worker_id", workerKey), zap.Error(err), zap.Duration("delay", delay))
						continue
					}
				}
			}
		}
	}()
}
