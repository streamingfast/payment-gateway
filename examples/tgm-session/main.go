package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/streamingfast/dauth"
	"github.com/streamingfast/dsession"
	"github.com/streamingfast/payment-gateway/auth"
	"github.com/streamingfast/payment-gateway/session"
	"go.uber.org/zap"
)

func init() {
	auth.Register()
	session.Register()
}

// this example demonstrates how to use the payment-gateway/session package to manage worker sessions
func main() {
	apiKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		fmt.Println("API_KEY environment variable is not set")
		return
	}

	// First, authenticate to get user context
	authPluginDSN := "tgm://"
	authenticator, err := dauth.New(authPluginDSN, zap.NewNop())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	headers := map[string][]string{
		"X-Api-Key": {apiKey},
	}

	ctx, err = authenticator.Authenticate(ctx, "/", headers, "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	// Now create session pool
	sessionPluginDSN := "tgm://session.thegraph.market?request-keep-alive-delay=30s&default-max-request-per-user=5"
	sessionPool, err := dsession.New(sessionPluginDSN, zap.NewNop())
	if err != nil {
		log.Fatal(err)
	}

	// Get session details from auth context
	trustedHeaders := dauth.FromContext(ctx)
	userID := trustedHeaders.UserID()
	apiKeyID := trustedHeaders.APIKeyID()

	// Borrow a session
	sessionKey, err := sessionPool.Get(ctx, "example-service", userID, apiKeyID, "trace-123", func(err error) {
		fmt.Printf("Session error: %v\n", err)
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Session borrowed successfully!\n")
	fmt.Printf("Session Key: %s\n", sessionKey)
	fmt.Printf("User ID: %s\n", userID)
	fmt.Printf("API Key ID: %s\n", apiKeyID)

	// Use the session for your work here...
	fmt.Println("Using session for work...")

	// Release the session when done
	sessionPool.Release(sessionKey, apiKeyID)
	fmt.Println("Session released successfully!")
}
