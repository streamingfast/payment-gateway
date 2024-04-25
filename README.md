## StreamingFast Payment Gateway

This repository contains public interfaces, clients and other utilities for the StreamingFast Payment Gateway, which is the all encompassing Gateway for consuming Substreams and Firehose services as well as tooling for indexer to report their usage.

### Usage Reporting

This library enables the `paymentGateway://...` metering plugin that can be hooked into your application for reporting usages.

To register the scheme and obtain a metering event emitter, you can use the following snippet:

> [!NOTE]
> Full example with import(s) at [./examples/payment-gateway-metering](./examples/payment-gateway-metering/main.go).

```go
func init() {
	metering.Register()
}

func main() {
	pluginDSN := os.ExpandEnv("paymentGateway://payment.gateway.streamingfast.io?network=eth-mainnet&token=${SF_API_TOKEN}")

	eventEmitter, err := dmetering.New(pluginDSN, zap.NewNop())
	if err != nil {
		panic(err)
	}

	// Start emitting events
}
```

> [!NOTE]
> Only selected indexer can try the public payment gateway for now.

#### gRPC

> [!NOTE]
> Full example with import(s) at [./examples/payment-gateway-usage-grpc](./examples/payment-gateway-usage-grpc/main.go).

```go
func main() {
	conn, err := dgrpc.NewClientConn("payment.gateway.streamingfast.io:443",
		dgrpc.WithMustAutoTransportCredentials(false, false, false),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})}),
	)
	failOnError(err)

	client := pbgateway.NewUsageServiceClient(conn)

	client := pbgateway.NewUsageServiceClient(conn)
	response, err := client.Report(context.TODO(), &pbgateway.ReportRequest{
		// Events: []*pbpbmetering.Event{ ... },
	})
	failOnError(err)

	// Don't forget to close the connection
	conn.Close()
}
```