package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/payment-gateway/auth"
	"go.uber.org/zap"
)

func init() {
	auth.Register()
}

// this example demonstrates how to use the payment-gateway/auth package to authenticate requests that use an API key
// will also work with an JWT passed in the Authorization header
func main() {
	apiKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		fmt.Println("API_KEY environment variable is not set")
		return
	}

	pluginDSN := "paymentgateway://"
	// pluginDSN := "paymentgateway://auth.thegraph.market?pubkeyurl=https://auth.thegraph.market/.well-known/jwks.json" // default values

	authenticator, err := dauth.New(pluginDSN, zap.NewNop())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	path := "/"
	_ = apiKey
	headers := map[string][]string{
		"X-Api-Key": {apiKey},
		//"authorization": {"Bearer " + jwt},
	}
	ipAddress := "127.0.0.123"

	ctx, err = authenticator.Authenticate(ctx, path, headers, ipAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Authentication successful!\nHeaders:")
	trustedHeaders := dauth.FromContext(ctx)
	for k, v := range trustedHeaders {
		fmt.Printf("  %s: %s\n", k, v)
	}

}
