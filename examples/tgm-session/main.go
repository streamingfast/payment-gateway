package main

import (
	"context"
	"fmt"
	"log"

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
	ctx := context.Background()

	// Create session pool
	sessionPluginDSN := "tgm://session.thegraph.market?request-keep-alive-delay=30s&default-max-request-per-user=5&indexer-api-key=server_1234567890abcdef"
	sessionPool, err := dsession.New(sessionPluginDSN, zap.NewNop())
	if err != nil {
		log.Fatal(err)
	}

	userID := "user123"
	apiKeyID := "apikey123"

	// Note: If you used dauth authenticator, you would normally get them from auth context, like this:
	//   trustedHeaders := dauth.FromContext(ctx)
	//   userID := trustedHeaders.UserID()
	//   apiKeyID := trustedHeaders.APIKeyID()

	// Borrow a session
	sessionKey, err := sessionPool.Get(ctx, "example-service", userID, apiKeyID, "trace-123", func(err error) {
		fmt.Printf("Session error: %v\n", err)
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Session borrowed successfully!\n")
	fmt.Printf("Session Key: %s\n", sessionKey)

	// Use the session for your work here...

	// Release the session when done
	sessionPool.Release(sessionKey)
	fmt.Println("Session released!")
}
