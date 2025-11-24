package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/streamingfast/dsession"
	pbworker "github.com/streamingfast/worker-pool-protocol/pb/sf/worker/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeWorkerPoolClient is a thread-safe fake implementing pbworker.WorkerPoolClient
type fakeWorkerPoolClient struct {
	mu sync.Mutex

	nextSessionID int
	nextWorkerID  int

	// configuration for behavior
	keepAliveErrOnce error // if set, first KeepAlive returns this error

	// recordings
	borrowRequests []*pbworker.BorrowWorkerRequest
	keepAlives     []*pbworker.KeepAliveRequest
	returned       []*pbworker.ReturnWorkerRequest
}

func newFakeWorkerPoolClient() *fakeWorkerPoolClient {
	return &fakeWorkerPoolClient{nextSessionID: 1, nextWorkerID: 1}
}

func (f *fakeWorkerPoolClient) BorrowWorker(ctx context.Context, in *pbworker.BorrowWorkerRequest, _ ...grpc.CallOption) (*pbworker.BorrowWorkerResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.borrowRequests = append(f.borrowRequests, in)

	// Distinguish session vs worker by MaxWorkerForTraceId field (0 for session per our usage)
	if in.GetMaxWorkerForTraceId() == 0 {
		key := fmt.Sprintf("session-%d", f.nextSessionID)
		f.nextSessionID++
		return &pbworker.BorrowWorkerResponse{
			Status:    pbworker.BorrowWorkerResponse_borrowed,
			WorkerKey: key,
		}, nil
	}

	// Worker borrow
	key := fmt.Sprintf("worker-%d", f.nextWorkerID)
	f.nextWorkerID++
	return &pbworker.BorrowWorkerResponse{
		Status:    pbworker.BorrowWorkerResponse_borrowed,
		WorkerKey: key,
	}, nil
}

// KeepAlive simulates a keep alive, optionally failing once with configured error
func (f *fakeWorkerPoolClient) KeepAlive(ctx context.Context, in *pbworker.KeepAliveRequest, _ ...grpc.CallOption) (*pbworker.KeepAliveResponse, error) {
	f.mu.Lock()
	f.keepAlives = append(f.keepAlives, in)
	var err error
	if f.keepAliveErrOnce != nil {
		err = f.keepAliveErrOnce
		f.keepAliveErrOnce = nil
	}
	f.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return &pbworker.KeepAliveResponse{}, nil
}

func (f *fakeWorkerPoolClient) ReturnWorker(ctx context.Context, in *pbworker.ReturnWorkerRequest, _ ...grpc.CallOption) (*pbworker.ReturnWorkerResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.returned = append(f.returned, in)
	return &pbworker.ReturnWorkerResponse{}, nil
}

// implement unused interface method to satisfy WorkerPoolClient
func (f *fakeWorkerPoolClient) WorkersState(ctx context.Context, in *pbworker.WorkersStateRequest, _ ...grpc.CallOption) (*pbworker.WorkersStateResponse, error) {
	return &pbworker.WorkersStateResponse{}, nil
}

// Ensure fake satisfies the interface at compile time
var _ pbworker.WorkerPoolClient = (*fakeWorkerPoolClient)(nil)

func TestSessionMutex_WithKeepAliveAndRelease_NoPanic(t *testing.T) {
	t.Parallel()

	cfg := &Config{RequestKeepAliveDelay: 10 * time.Millisecond, MinimalWorkerLifeDuration: 10 * time.Millisecond}
	fake := newFakeWorkerPoolClient()
	pool := &tgmSessionPool{
		config:                 cfg,
		logger:                 zap.NewNop(),
		remoteWorkerPoolClient: fake,
		sessions:               make(map[string]*sessionInfo),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Borrow a session (worker) -> starts keepalive
	key, err := pool.Get(ctx, "svc", "org1", "api1", "trace-1", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Allow a couple of keepalive ticks
	time.Sleep(25 * time.Millisecond)

	// Concurrently release while keepalive may be running
	done := make(chan struct{})
	go func() {
		pool.Release(key)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("Release did not complete in time (possible deadlock)")
	}
}

func TestGetWorker_Tracking_AddRemove_UnderMutex(t *testing.T) {
	t.Parallel()

	cfg := &Config{RequestKeepAliveDelay: time.Hour, MinimalWorkerLifeDuration: 10 * time.Millisecond}
	fake := newFakeWorkerPoolClient()
	pool := &tgmSessionPool{
		config:                 cfg,
		logger:                 zap.NewNop(),
		remoteWorkerPoolClient: fake,
		sessions:               make(map[string]*sessionInfo),
	}

	ctx := context.Background()
	sessionKey, err := pool.Get(ctx, "svc", "org2", "api2", "trace-2", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Borrow several workers concurrently
	const n = 5
	var wg sync.WaitGroup
	workerKeys := make([]string, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			wkey, err := pool.GetWorker(ctx, "svc", sessionKey, n)
			workerKeys[idx] = wkey
			errs[idx] = err
		}(i)
	}
	wg.Wait()

	for i, e := range errs {
		if e != nil {
			t.Fatalf("GetWorker %d failed: %v", i, e)
		}
		if workerKeys[i] == "" {
			t.Fatalf("GetWorker %d returned empty key", i)
		}
	}

	// Release some workers explicitly
	for i := 0; i < n/2; i++ {
		pool.ReleaseWorker(workerKeys[i])
	}

	// Release the session, which should release remaining workers
	pool.Release(sessionKey)

	// Wait until the fake observed returns for all workers and the session
	deadline := time.Now().Add(2 * time.Second)
	for {
		fake.mu.Lock()
		count := len(fake.returned)
		fake.mu.Unlock()
		if count >= n+1 { // n workers + 1 session
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for ReturnWorker calls, have %d want %d", count, n+1)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestKeepAlive_OnPermissionDenied_TriggersOnError(t *testing.T) {
	t.Parallel()

	cfg := &Config{RequestKeepAliveDelay: 5 * time.Millisecond, MinimalWorkerLifeDuration: 10 * time.Millisecond}
	fake := newFakeWorkerPoolClient()
	fake.keepAliveErrOnce = status.Error(codes.PermissionDenied, "nope")

	pool := &tgmSessionPool{
		config:                 cfg,
		logger:                 zap.NewNop(),
		remoteWorkerPoolClient: fake,
		sessions:               make(map[string]*sessionInfo),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	onError := func(err error) {
		errCh <- err
	}

	_, err := pool.Get(ctx, "svc", "org3", "api3", "trace-3", onError)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	select {
	case e := <-errCh:
		if !errors.Is(e, dsession.ErrPermissionDenied) {
			t.Fatalf("expected ErrPermissionDenied, got %v", e)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("onError was not called in time")
	}
}
