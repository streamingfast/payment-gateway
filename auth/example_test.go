package auth_test

import (
	"context"
	"fmt"
	"log"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/payment-gateway/auth"
	"go.uber.org/zap"
)

func Example() {
	// Register the payment gateway authenticator
	auth.Register()

	// Example configuration string
	// This would typically come from environment variables or config files
	config := "paymentgateway://auth.thegraph.market?pubkeyurl=https://auth.thegraph.market/.well-known/jwks.json&reissue-jwt-max-age-secs=600&key=server_key"

	// Create a logger
	logger := zap.NewNop()

	// Create the authenticator using dauth
	authenticator, err := dauth.New(config, logger)
	if err != nil {
		log.Fatalf("Failed to create authenticator: %v", err)
	}

	// Example 1: Authenticate with JWT token
	ctx := context.Background()
	headers := map[string][]string{
		"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
	}
	ipAddress := "192.168.1.100"

	authenticatedCtx, err := authenticator.Authenticate(ctx, "/api/v1/query", headers, ipAddress)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return
	}

	// Extract trusted headers from context
	trustedHeaders := dauth.FromContext(authenticatedCtx)
	fmt.Printf("User ID: %s\n", trustedHeaders[dauth.SFHeaderUserID])
	fmt.Printf("API Key ID: %s\n", trustedHeaders[dauth.SFHeaderApiKeyID])
	fmt.Printf("IP Address: %s\n", trustedHeaders[dauth.SFHeaderIP])

	// Example 2: Authenticate with API key
	headers = map[string][]string{
		"X-API-Key": {"sk_live_1234567890abcdef"},
	}

	authenticatedCtx, err = authenticator.Authenticate(ctx, "/api/v1/query", headers, ipAddress)
	if err != nil {
		fmt.Printf("Authentication with API key failed: %v\n", err)
		return
	}

	// The API key will be exchanged for a JWT token internally
	trustedHeaders = dauth.FromContext(authenticatedCtx)
	fmt.Printf("Authenticated via API key - User ID: %s\n", trustedHeaders[dauth.SFHeaderUserID])

	// Example 3: Access feature configurations from JWT
	if parallelJobs, ok := trustedHeaders["x-sf-substreams-parallel-jobs"]; ok {
		fmt.Printf("Substreams parallel jobs: %s\n", parallelJobs)
	}

	// Output:
	// User ID: user-123
	// API Key ID: api-key-456
	// IP Address: 192.168.1.100
	// Authenticated via API key - User ID: user-789
	// Substreams parallel jobs: 10
}

func ExampleAuthenticator_configuration() {
	// Configuration examples for different scenarios

	// 1. Production configuration with JWK URL
	prodConfig := "paymentgateway://auth.production.com?" +
		"pubkeyurl=https://auth.production.com/.well-known/jwks.json&" +
		"reissue-jwt-max-age-secs=300&" +
		"key=server_production_key"

	// 2. Development configuration with plaintext HTTP
	devConfig := "paymentgateway://localhost:8080?" +
		"plaintext=true&" +
		"insecure=true&" +
		"pubkeyurl=http://localhost:8080/.well-known/jwks.json&" +
		"reissue-jwt-max-age-secs=60"

	// 3. Configuration with base64-encoded public key (for offline verification)
	offlineConfig := "paymentgateway://auth.example.com?" +
		"pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0K...&" +
		"reissue-jwt-max-age-secs=600"

	fmt.Printf("Production: %s\n", prodConfig)
	fmt.Printf("Development: %s\n", devConfig)
	fmt.Printf("Offline: %s\n", offlineConfig)

	// Output:
	// Production: paymentgateway://auth.production.com?pubkeyurl=https://auth.production.com/.well-known/jwks.json&reissue-jwt-max-age-secs=300&key=server_production_key
	// Development: paymentgateway://localhost:8080?plaintext=true&insecure=true&pubkeyurl=http://localhost:8080/.well-known/jwks.json&reissue-jwt-max-age-secs=60
	// Offline: paymentgateway://auth.example.com?pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0K...&reissue-jwt-max-age-secs=600
}
