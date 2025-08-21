#  The Graph Market plugins for authentication and metering

This repository contains public interfaces, clients and other utilities for The Graph market authentication and metering, used for consuming Substreams, Firehose and tokenAPI services as well as tooling for indexer to report their usage.

It is intended for developers that would like to offer their services on this platform.

### Authentication plugin

* Doc: [./auth/README.md](./auth/README.md)
* Example: [./examples/payment-gateway-authentication](./examples/payment-gateway-authentication/main.go)

### Metering plugin

This library enables the `tgm://...` (formerly `paymentGateway://...`) metering plugin that can be hooked into your application for reporting usages.

To register the scheme and obtain a metering event emitter, you can use the following snippet:

> [!NOTE]
> Full example with import(s) and code annotations at [./examples/tgm-metering](./examples/payment-gateway-metering/main.go).

```go
func init() {
	// Register tgm:// as a valid metering plugin, refers to Register documentation for extra details
	metering.Register()
}

func main() {
	pluginDSN := "tgm://abp.thegraph.market?network=eth-mainnet&token=${API_TOKEN}"

	eventEmitter, err := dmetering.New(pluginDSN, zap.NewNop())
	if err != nil {
		panic(err)
	}

	// Start emitting events
}
```

### gRPC

> [!NOTE]
> Full example with import(s) and code annotations at [./examples/tgm-usage-grpc](./examples/payment-gateway-usage-grpc/main.go).

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
