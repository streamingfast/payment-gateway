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

type sessionInfo struct {
	userID   string
	apiKeyID string
	traceID  string
	workers  map[string]struct{} // Track worker keys for this session
	closer   chan struct{}       // Channel to signal session closure
}

type tgmSessionPool struct {
	config                 *Config
	logger                 *zap.Logger
	remoteWorkerPoolClient pbworker.WorkerPoolClient
	conn                   *grpc.ClientConn
	sessions               map[string]*sessionInfo // Map sessionKey -> session info
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
		sessions:               make(map[string]*sessionInfo),
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

	details := ""
	if maxWorkers := resp.WorkerState.GetMaxWorkers(); maxWorkers != 0 {
		details = fmt.Sprintf(" (active sessions: %d/%d)", maxWorkers, maxWorkers)
	}

	if workerStatus == pbworker.BorrowWorkerResponse_resource_exhausted {
		t.logger.Info("worker pool is exhausted", zap.String("worker_key", key), zap.String("status", workerStatus.String()))
		return "", fmt.Errorf("%w%s", dsession.ErrConcurrentStreamLimitExceeded, details)
	}

	// Start keep-alive for borrowed workers
	if workerStatus == pbworker.BorrowWorkerResponse_borrowed {

		done := make(chan struct{})
		t.sessionMutex.Lock()
		// Store session info for worker management
		t.sessions[key] = &sessionInfo{
			userID:   userID,
			apiKeyID: apiKeyID,
			traceID:  traceID,
			workers:  make(map[string]struct{}),
			closer:   done,
		}
		t.sessionMutex.Unlock()

		t.startKeepAlive(ctx, done, key, onError)
	}

	t.logger.Debug("borrowed request worker", zap.String("worker_key", key))

	return key, nil
}

func (t *tgmSessionPool) Release(sessionKey string) {
	go func() {
		t.sessionMutex.Lock()
		sessionInfo := t.sessions[sessionKey]

		// Collect all workers to release and close the session
		var workersToRelease []string
		var done chan struct{}
		if sessionInfo != nil {
			for workerKey := range sessionInfo.workers {
				workersToRelease = append(workersToRelease, workerKey)
			}
			done = sessionInfo.closer
			delete(t.sessions, sessionKey)
		}
		t.sessionMutex.Unlock()

		// Close the done channel after releasing the lock
		if done != nil {
			close(done)
		}

		// Release all workers associated with this session
		for _, workerKey := range workersToRelease {
			t.releaseWorkerInternal(workerKey)
		}

		resp, err := t.remoteWorkerPoolClient.ReturnWorker(context.Background(),
			&pbworker.ReturnWorkerRequest{
				WorkerKey: sessionKey,
			},
			grpc.WaitForReady(false),
		)

		t.logger.Debug("returned request worker", zap.String("key", sessionKey), zap.Stringer("status", resp.GetStatus()), zap.Error(err))
	}()
}

func (t *tgmSessionPool) GetWorker(ctx context.Context, serviceName string, sessionKey string, maxWorkersPerSession int) (string, error) {
	// Look up session info
	t.sessionMutex.Lock()
	sessionInfo := t.sessions[sessionKey]
	if sessionInfo == nil {
		t.sessionMutex.Unlock()
		return "", fmt.Errorf("%w: session key %s not found", dsession.ErrSessionNotFound, sessionKey)
	}
	// Copy the values we need while holding the lock
	userID := sessionInfo.userID
	apiKeyID := sessionInfo.apiKeyID
	traceID := sessionInfo.traceID
	t.sessionMutex.Unlock()

	resp, err := t.remoteWorkerPoolClient.BorrowWorker(ctx,
		&pbworker.BorrowWorkerRequest{
			Service:             serviceName,
			UserId:              userID,
			ApiKeyId:            apiKeyID,
			TraceId:             traceID,
			MaxWorkerForTraceId: int64(maxWorkersPerSession),
		},
		grpc.WaitForReady(false),
	)

	if err != nil {
		// Map gRPC errors to dsession errors
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return "", fmt.Errorf("%w: session not found", dsession.ErrSessionNotFound)
			case codes.ResourceExhausted:
				return "", fmt.Errorf("%w: maximum workers per session exceeded", dsession.ErrWorkersLimitExceeded)
			}
		}
		return "", fmt.Errorf("failed to borrow worker: %w", err)
	}

	workerKey := resp.WorkerKey
	workerStatus := resp.Status

	details := ""
	if maxWorkers := resp.WorkerState.GetMaxWorkers(); maxWorkers != 0 {
		details = fmt.Sprintf(" (active workers: %d/%d)", maxWorkers, maxWorkers)
	}

	if workerStatus == pbworker.BorrowWorkerResponse_resource_exhausted {
		t.logger.Info("worker limit exceeded", zap.String("worker_key", workerKey), zap.String("status", workerStatus.String()))
		return "", fmt.Errorf("%w%s", dsession.ErrWorkersLimitExceeded, details)
	}

	// Track this worker under the session
	t.sessionMutex.Lock()
	sessionInfo = t.sessions[sessionKey]
	if sessionInfo == nil {
		t.sessionMutex.Unlock()
		// Session was released, immediately release the newly acquired worker
		go t.releaseWorkerInternal(workerKey)
		return "", fmt.Errorf("%w: session key %s was released", dsession.ErrSessionNotFound, sessionKey)
	}
	sessionInfo.workers[workerKey] = struct{}{}
	t.sessionMutex.Unlock()

	t.logger.Debug("borrowed worker", zap.String("worker_key", workerKey), zap.String("session_key", sessionKey), zap.Int("max_workers", maxWorkersPerSession))

	return workerKey, nil
}

func (t *tgmSessionPool) ReleaseWorker(workerKey string) {
	// Remove worker from session tracking
	t.sessionMutex.Lock()
	for _, sessionInfo := range t.sessions {
		delete(sessionInfo.workers, workerKey)
	}
	t.sessionMutex.Unlock()

	// Release worker in a goroutine (fire-and-forget)
	go t.releaseWorkerInternal(workerKey)
}

func (t *tgmSessionPool) releaseWorkerInternal(workerKey string) {
	resp, err := t.remoteWorkerPoolClient.ReturnWorker(context.Background(),
		&pbworker.ReturnWorkerRequest{
			WorkerKey: workerKey,
		},
		grpc.WaitForReady(false),
	)

	t.logger.Debug("returned worker", zap.String("key", workerKey), zap.Stringer("status", resp.GetStatus()), zap.Error(err))
}

// createApiKeyInterceptor creates a gRPC unary interceptor that adds the X-Api-Key header
func createApiKeyInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add the X-Api-Key header to the outgoing metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "X-Api-Key", apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// startKeepAlive starts the keep-alive goroutine for a borrowed session and its workers
func (t *tgmSessionPool) startKeepAlive(ctx context.Context, done <-chan struct{}, sessionKey string, onError func(error)) {
	go func() {
		// Use a ticker for consistent intervals regardless of operation duration
		// The ticker will fire at regular intervals from when it starts, not from when each tick is consumed
		ticker := time.NewTicker(t.config.RequestKeepAliveDelay)
		defer ticker.Stop()

		// Track if we're in error recovery mode with 1-second intervals
		errorMode := false

		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-ticker.C:
				// The ticker ensures consistent intervals - it ticks at regular intervals
				// regardless of how long the keep-alive operations take

				// Get session info and workers to keep alive
				t.sessionMutex.Lock()
				sessionInfo := t.sessions[sessionKey]
				if sessionInfo == nil {
					t.sessionMutex.Unlock()
					return // Session was released
				}
				apiKeyID := sessionInfo.apiKeyID
				workerKeys := make([]string, 0, len(sessionInfo.workers))
				for workerKey := range sessionInfo.workers {
					workerKeys = append(workerKeys, workerKey)
				}
				t.sessionMutex.Unlock()

				hadError := false

				// Keep session alive
				_, err := t.remoteWorkerPoolClient.KeepAlive(
					ctx,
					&pbworker.KeepAliveRequest{
						WorkerKey: sessionKey,
						ApiKeyId:  apiKeyID,
					},
					grpc.WaitForReady(false),
				)
				if err != nil {
					hadError = true
					t.logger.Error("failed to call keep session alive", zap.String("session_key", sessionKey), zap.Error(err))
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
					}
				}

				// Keep workers alive
				for _, workerKey := range workerKeys {
					_, err := t.remoteWorkerPoolClient.KeepAlive(
						ctx,
						&pbworker.KeepAliveRequest{
							WorkerKey: workerKey,
							ApiKeyId:  apiKeyID,
						},
						grpc.WaitForReady(false),
					)
					if err != nil {
						hadError = true
						t.logger.Error("failed to call keep worker alive", zap.String("worker_key", workerKey), zap.Error(err))
					}
				}

				// On error, switch to 1-second retry interval
				// On success after error, switch back to normal interval
				if hadError && !errorMode {
					ticker.Reset(time.Second)
					errorMode = true
					t.logger.Info("switched to error recovery mode with 1 second interval", zap.String("session_key", sessionKey))
				} else if !hadError && errorMode {
					ticker.Reset(t.config.RequestKeepAliveDelay)
					errorMode = false
					t.logger.Info("recovered from error, switched back to normal interval", zap.String("session_key", sessionKey))
				}
			}
		}
	}()
}
