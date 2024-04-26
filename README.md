## StreamingFast Payment Gateway

This repository contains public interfaces, clients and other utilities for the StreamingFast Payment Gateway, which is the all encompassing Gateway for consuming Substreams and Firehose services as well as tooling for indexer to report their usage.

### Service Developers

This section serves as the documentation for developers that would like to offer their services on StreamingFast Payment Gateway. If you are **only** a consumer of services, this section is not for you.

### Usage Reporting

This library enables the `paymentGateway://...` metering plugin that can be hooked into your application for reporting usages.

To register the scheme and obtain a metering event emitter, you can use the following snippet:

> [!NOTE]
> Full example with import(s) and code annotations at [./examples/payment-gateway-metering](./examples/payment-gateway-metering/main.go).

```go
func init() {
	// Register paymentGateway:// as a valid metering plugin, refers to Register documentation for extra details
	metering.Register()
}

func main() {
	pluginDSN := "paymentGateway://abp.thegraph.market?network=eth-mainnet&token=${API_TOKEN}"

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
> Full example with import(s) and code annotations at [./examples/payment-gateway-usage-grpc](./examples/payment-gateway-usage-grpc/main.go).

```go
func main() {
	token := os.Getenv("API_TOKEN")

	conn, err := dgrpc.NewClientConn("abp.thegraph.market:443",
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