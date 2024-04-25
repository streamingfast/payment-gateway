package main

import (
	"context"
	"fmt"
	"os"

	"github.com/streamingfast/cli"
	"github.com/streamingfast/dgrpc"
	pbgateway "github.com/streamingfast/payment-gateway/pb/sf/gateway/payment/v1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func main() {
	endpoint := osEnvOr("PAYMENT_GATEWAY_ENDPOINT", "payment.gateway.streamingfast.io:443")
	insecure := osEnvOr("PAYMENT_GATEWAY_INSECURE", "false")
	plainText := osEnvOr("PAYMENT_GATEWAY_PLAINTEXT", "false")
	token := osEnv("SF_API_TOKEN")

	conn, err := dgrpc.NewClientConn(endpoint,
		dgrpc.WithMustAutoTransportCredentials(insecure == "true", plainText == "true", false),
		grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})}),
	)
	cli.NoError(err, "unable to create external gRPC client")

	client := pbgateway.NewUsageServiceClient(conn)
	response, err := client.Report(context.TODO(), &pbgateway.ReportRequest{
		// Events: []*pbpbmetering.Event{ ... },
	})
	cli.NoError(err, "unable to report events")

	// Do something with the response
	_ = response

	// Don't forget to close the connection
	conn.Close()
}

func osEnv(key string) string {
	if value, exists := os.LookupEnv(key); !exists {
		panic(fmt.Errorf("missing required environment variable %q", key))
	} else {
		return value
	}
}

func osEnvOr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
