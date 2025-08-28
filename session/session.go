package session

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dsession"
	pbworker "github.com/streamingfast/worker-pool-protocol/pb/sf/worker/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
	globalRequestPool      *GlobalRequestPool
	remoteWorkerPoolClient pbworker.WorkerPoolClient
	conn                   *grpc.ClientConn
}

func newTGMSessionPool(config *Config, logger *zap.Logger) (dsession.SessionPool, error) {
	logger = logger.Named("tgm-session-pool")

	// Determine the connection scheme (removed unused variable)

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

	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session endpoint %s: %w", config.Endpoint, err)
	}

	client := pbworker.NewWorkerPoolClient(conn)

	globalPool := NewGlobalRequestPool(
		client,
		config.RequestKeepAliveDelay,
		config.DefaultMaxRequestPerUser,
		config.DefaultMinimalWorkerLifeDuration,
		logger,
	)

	pool := &tgmSessionPool{
		config:                 config,
		logger:                 logger,
		globalRequestPool:      globalPool,
		remoteWorkerPoolClient: client,
		conn:                   conn,
	}

	return pool, nil
}

func (t *tgmSessionPool) Get(ctx context.Context, serviceName string, userID string, apiKeyID string, traceID string, onError func(error)) (string, error) {
	borrowedRequest, err := t.globalRequestPool.BorrowRequest(ctx, serviceName, userID, apiKeyID, traceID)
	if err != nil {
		return "", fmt.Errorf("failed to borrow session: %w", err)
	}

	// Store the borrowed request for later release
	t.globalRequestPool.storeBorrowedRequest(borrowedRequest.key, borrowedRequest)

	return borrowedRequest.key, nil
}

func (t *tgmSessionPool) Release(sessionKey, apiKeyID string) {
	if borrowedRequest := t.globalRequestPool.getBorrowedRequest(sessionKey); borrowedRequest != nil {
		t.globalRequestPool.ReturnRequest(borrowedRequest)
		t.globalRequestPool.removeBorrowedRequest(sessionKey)
	}
}

// tgmSession implements Session
type tgmSession struct {
	borrowedRequest   *BorrowedRequest
	globalRequestPool *GlobalRequestPool
	logger            *zap.Logger
}

func (s *tgmSession) WorkerKey() string {
	return s.borrowedRequest.key
}

func (s *tgmSession) Status() string {
	return s.borrowedRequest.status.String()
}

func (s *tgmSession) IsResourceExhausted() bool {
	return s.borrowedRequest.status == pbworker.BorrowWorkerResponse_resource_exhausted
}

func (s *tgmSession) WorkerState() *pbworker.WorkersState {
	return s.borrowedRequest.state
}

func (s *tgmSession) Close() error {
	s.globalRequestPool.ReturnRequest(s.borrowedRequest)
	return nil
}

// GlobalRequestPool manages worker requests
type GlobalRequestPool struct {
	userBorrowedRequest              map[string]uint64
	borrowedRequests                 map[string]*BorrowedRequest
	borrowedRequestMutex             sync.Mutex
	remoteWorkerPoolClient           pbworker.WorkerPoolClient
	requestKeepAliveDelay            time.Duration
	logger                           *zap.Logger
	defaultMaxRequestPerUser         uint64
	defaultMinimalWorkerLifeDuration time.Duration
}

func NewGlobalRequestPool(remoteWorkerPoolClient pbworker.WorkerPoolClient, requestKeepAliveDelay time.Duration, defaultMaxRequestPerUser uint64, defaultMinimalWorkerLifeDuration time.Duration, logger *zap.Logger) *GlobalRequestPool {
	logger = logger.Named("global-request-pool")

	return &GlobalRequestPool{
		userBorrowedRequest:              make(map[string]uint64),
		borrowedRequests:                 make(map[string]*BorrowedRequest),
		borrowedRequestMutex:             sync.Mutex{},
		remoteWorkerPoolClient:           remoteWorkerPoolClient,
		requestKeepAliveDelay:            requestKeepAliveDelay,
		defaultMaxRequestPerUser:         defaultMaxRequestPerUser,
		defaultMinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
		logger:                           logger,
	}
}

func (p *GlobalRequestPool) BorrowRequest(ctx context.Context, serviceName string, userID string, apiKeyID string, traceID string) (*BorrowedRequest, error) {
	resp, err := p.remoteWorkerPoolClient.BorrowWorker(ctx,
		&pbworker.BorrowWorkerRequest{
			Service:  serviceName,
			UserId:   userID,
			ApiKeyId: apiKeyID,
			TraceId:  traceID,
		},
		grpc.WaitForReady(false),
	)

	if err != nil {
		return nil, err
	}

	key := resp.WorkerKey
	workerStatus := resp.Status
	state := resp.WorkerState
	minimalWorkerLifeDuration := resp.MinimalWorkerLifeDuration.AsDuration()

	r := NewBorrowedRequest(key, userID, workerStatus, state, minimalWorkerLifeDuration, p.logger)

	if workerStatus == pbworker.BorrowWorkerResponse_resource_exhausted {
		p.logger.Info("worker pool is exhausted", zap.String("worker_key", key), zap.String("status", workerStatus.String()))
		return r, nil
	}

	p.borrowedRequestMutex.Lock()
	p.userBorrowedRequest[userID]++
	p.borrowedRequestMutex.Unlock()

	r.startKeepAlive(ctx, p.requestKeepAliveDelay, p.remoteWorkerPoolClient)

	p.logger.Info("borrowed request worker", zap.String("worker_key", key))

	return r, nil
}

func (p *GlobalRequestPool) storeBorrowedRequest(key string, request *BorrowedRequest) {
	p.borrowedRequestMutex.Lock()
	defer p.borrowedRequestMutex.Unlock()
	p.borrowedRequests[key] = request
}

func (p *GlobalRequestPool) getBorrowedRequest(key string) *BorrowedRequest {
	p.borrowedRequestMutex.Lock()
	defer p.borrowedRequestMutex.Unlock()
	return p.borrowedRequests[key]
}

func (p *GlobalRequestPool) removeBorrowedRequest(key string) {
	p.borrowedRequestMutex.Lock()
	defer p.borrowedRequestMutex.Unlock()
	delete(p.borrowedRequests, key)
}

func (p *GlobalRequestPool) ReturnRequest(r *BorrowedRequest) {
	r.StopKeepAlive()

	p.borrowedRequestMutex.Lock()
	p.userBorrowedRequest[r.userID]--
	p.borrowedRequestMutex.Unlock()

	resp, err := p.remoteWorkerPoolClient.ReturnWorker(context.Background(),
		&pbworker.ReturnWorkerRequest{
			WorkerKey:                 r.key,
			MinimalWorkerLifeDuration: durationpb.New(r.minimalWorkerLifeDuration),
		},
		grpc.WaitForReady(false),
	)

	if err != nil {
		p.logger.Error("returning request worker", zap.Error(err))
	} else {
		p.logger.Info("returned request worker", zap.String("key", r.key), zap.Stringer("status", resp.Status))
	}
}

// BorrowedRequest represents a borrowed worker session
type BorrowedRequest struct {
	key                       string
	userID                    string
	status                    pbworker.BorrowWorkerResponse_BorrowStatus
	state                     *pbworker.WorkersState
	minimalWorkerLifeDuration time.Duration
	logger                    *zap.Logger
	done                      chan struct{}
}

func NewBorrowedRequest(key string, userID string, status pbworker.BorrowWorkerResponse_BorrowStatus, state *pbworker.WorkersState, minimalWorkerLifeDuration time.Duration, logger *zap.Logger) *BorrowedRequest {
	logger = logger.Named("borrowed-request")
	return &BorrowedRequest{
		key:                       key,
		userID:                    userID,
		status:                    status,
		state:                     state,
		minimalWorkerLifeDuration: minimalWorkerLifeDuration,
		done:                      make(chan struct{}),
		logger:                    logger,
	}
}

func (r *BorrowedRequest) startKeepAlive(ctx context.Context, delay time.Duration, remoteWorkerPoolClient pbworker.WorkerPoolClient) {
	apiKeyID := dauth.FromContext(ctx).APIKeyID()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.done:
				return
			case <-time.After(delay):
				_, err := remoteWorkerPoolClient.KeepAlive(
					ctx,
					&pbworker.KeepAliveRequest{
						WorkerKey: r.key,
						ApiKeyId:  apiKeyID,
					},
					grpc.WaitForReady(false),
				)
				if err != nil {
					// FIXME: add some retry with smaller delay, handle specific errors here
					r.logger.Error("failed to call keep request worker alive", zap.String("worker_id", r.key), zap.Error(err))
					//reqctx.CancelFunc(ctx)(err)
					return
				}
			}
		}
	}()
}

func (r *BorrowedRequest) StopKeepAlive() {
	close(r.done)
}

// Helper function to extract trace ID from context
// This is a placeholder - you'll need to implement based on your tracing setup
func getTraceIDFromContext(ctx context.Context) string {
	// Implementation depends on your tracing framework (e.g., OpenTelemetry, Jaeger, etc.)
	// For now, return empty string
	return ""
}
