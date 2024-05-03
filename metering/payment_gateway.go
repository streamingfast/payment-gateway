package metering

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/dgrpc"
	"github.com/streamingfast/dmetering"
	pbmetering "github.com/streamingfast/dmetering/pb/sf/metering/v1"
	"github.com/streamingfast/dmetrics"
	pbgateway "github.com/streamingfast/payment-gateway/pb/sf/gateway/payment/v1"
	"github.com/streamingfast/shutter"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
)

// FIXME: This is super similar to https://github.com/dfuse-io/dmetering/blob/4d001fb2b225632174a30f7fb97e98c0d7fb904f/grpc/emitter.go.
// This is because this implementation does batching based on network. We should probably merge the two implementations
// and make the final dispath configurable to accomodate both payment gateway and metering services.

// Register registers the payment gateway emitter inside [dmetering] package by calling
// [dmetering.Register] with the proper factory function. The emitter accepts a URL with the following format:
//
//	paymentGateway://<endpoint>?network=<network>&token=${SF_API_TOKEN}[&insecure=true|false][&plaintext=true|false][&delay=1s][&buffer_size=1000][&panic_on_drop=true|false]
//
// The metering plugin will emit events to the payment gateway service pointed to by <endpoint> (required). The endpoint
// can contains a `:<port>` suffix to specify which port to use. If the port is not provided, 443 is assumed.
//
// The connection is secured by TLS by default. If the `insecure` query parameter is set to `true`, the connection will be made
// with TLS but without verifying the certificate. The `plaintext` is reserved for development purposes and the connection is
// made without TLS, you cannot use that on production endpoints since they require a <token> and that a <token> can be sent
// onlt if the connection is secured with TLS.
func Register() {
	dmetering.Register("paymentgateway", func(config string, logger *zap.Logger) (dmetering.EventEmitter, error) {
		configExpanded := os.ExpandEnv(config)

		c, err := newConfig(configExpanded)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config string %s: %w", config, err)
		}
		return newEmitter(c, logger)
	})
}

type CloseFunc func() error

type emitter struct {
	*shutter.Shutter
	config *Config

	activeBatch     []*pbmetering.Event
	buffer          chan dmetering.Event
	client          pbgateway.UsageServiceClient
	clientCloseFunc CloseFunc
	done            chan bool

	logger *zap.Logger
}

func newEmitter(config *Config, logger *zap.Logger) (dmetering.EventEmitter, error) {
	client, closeFunc, err := newGatewayClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create external gRPC client %w", err)
	}

	return newWithClient(config, client, closeFunc, logger)
}

func newWithClient(
	config *Config,
	client pbgateway.UsageServiceClient,
	closeFunc CloseFunc,
	logger *zap.Logger,
) (dmetering.EventEmitter, error) {
	e := &emitter{
		Shutter:         shutter.New(),
		config:          config,
		client:          client,
		clientCloseFunc: closeFunc,
		buffer:          make(chan dmetering.Event, config.BufferSize),
		activeBatch:     []*pbmetering.Event{},
		done:            make(chan bool, 1),
		logger:          logger.Named("metrics.emitter"),
	}

	dmetrics.Register(MetricSet)

	e.OnTerminating(func(err error) {
		e.logger.Info("received shutdown signal, waiting for launch loop to end", zap.Error(err))
		<-e.done
		e.flushAndCloseEvent()
		closeErr := e.clientCloseFunc()
		if closeErr != nil {
			e.logger.Warn("failed to close grpc client", zap.Error(closeErr))
		}
	})
	go e.launch()

	return e, nil
}

func (e *emitter) launch() {
	ticker := time.NewTicker(e.config.Delay)
	for {
		select {
		case <-e.Terminating():
			e.done <- true
			return
		case <-ticker.C:
			//e.logger.Debug("emitting events after ticker delay", zap.Int("count", len(e.activeBatch)))
			e.emit(e.activeBatch)
			e.activeBatch = []*pbmetering.Event{}
		case ev := <-e.buffer:
			e.activeBatch = append(e.activeBatch, ev.ToProto(e.config.Network))
		}
	}
}

func (e *emitter) flushAndCloseEvent() {
	close(e.buffer)

	t0 := time.Now()
	e.logger.Info("waiting for event flush to complete", zap.Int("count", len(e.buffer)))
	defer func() {
		e.logger.Info("event flushed", zap.Duration("elapsed", time.Since(t0)))
	}()

	for {
		ev, ok := <-e.buffer
		protoEv := ev.ToProto(e.config.Network)
		if !ok {
			e.logger.Info("sending last events", zap.Int("count", len(e.activeBatch)))
			e.emit(e.activeBatch)
			return
		}
		e.activeBatch = append(e.activeBatch, protoEv)
	}
}

func (e *emitter) Emit(_ context.Context, ev dmetering.Event) {
	if ev.Endpoint == "" {
		e.logger.Warn("events must contain endpoint, dropping event", zap.Object("event", ev))
		return
	}

	if e.IsTerminating() {
		e.logger.Warn("emitter is shutting down cannot track event", zap.Object("event", ev))
		return
	}

	select {
	case e.buffer <- ev:
	default:
		if e.config.PanicOnDrop {
			panic(fmt.Errorf("failed to queue metric channel is full"))
		}
		DroppedEventCounter.Inc()
	}
}

func (e *emitter) emit(events []*pbmetering.Event) {
	if len(events) == 0 {
		return
	}
	e.logger.Debug("tracking events", zap.Int("count", len(events)))

	err := derr.RetryContext(context.Background(), 3, func(ctx context.Context) error {
		_, err := e.client.Report(context.Background(), &pbgateway.ReportRequest{Events: events})
		if err != nil {
			if dgrpc.IsGRPCErrorCode(err, codes.Unauthenticated) || dgrpc.IsGRPCErrorCode(err, codes.PermissionDenied) {
				return derr.NewFatalError(err)
			}

			MeteringGRPCRetryCounter.Inc()
			return err
		}

		return nil
	})

	if err != nil {
		MeteringGRPCErrCounter.Inc()
		e.logger.Warn("failed to emit event", zap.Error(err))
	}
}

func newGatewayClient(config *Config) (pbgateway.UsageServiceClient, CloseFunc, error) {
	opts := []grpc.DialOption{
		dgrpc.WithMustAutoTransportCredentials(config.Insecure, config.Plaintext, false),
	}
	if config.Token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: config.Token,
		})}))
	}

	conn, err := dgrpc.NewClientConn(config.Endpoint, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create external gRPC client: %w", err)
	}

	return pbgateway.NewUsageServiceClient(conn), conn.Close, nil
}
