package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streamingfast/dsession"
	"github.com/streamingfast/payment-gateway/session"
	"go.uber.org/zap"
)

func main() {
	// Register the TGM session pool
	session.Register()

	ctx := context.Background()

	// Create session pool
	sessionPluginDSN := "tgm://session.thegraph.market?request-keep-alive-delay=30s&default-max-request-per-user=5&indexer-api-key=server_1234567890abcdef"
	sessionPool, err := dsession.New(sessionPluginDSN, zap.NewNop())
	if err != nil {
		log.Fatal(err)
	}

	// User and service information
	userID := "user123"
	apiKeyID := "key456"
	traceID := "trace789"
	serviceName := "my-service"

	// Error handler for session errors
	onError := func(err error) {
		fmt.Printf("Session error: %v\n", err)
	}

	// Step 1: Get a session
	fmt.Println("Getting session...")
	sessionKey, err := sessionPool.Get(ctx, serviceName, userID, apiKeyID, traceID, onError)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
	}
	fmt.Printf("Got session: %s\n", sessionKey)

	// Step 2: Get workers under this session
	fmt.Println("\nGetting workers for the session...")

	// Get first worker
	worker1, err := sessionPool.GetWorker(ctx, serviceName, sessionKey, 3) // max 3 workers per session
	if err != nil {
		log.Fatalf("Failed to get worker 1: %v", err)
	}
	fmt.Printf("Got worker 1: %s\n", worker1)

	// Get second worker
	worker2, err := sessionPool.GetWorker(ctx, serviceName, sessionKey, 3)
	if err != nil {
		log.Fatalf("Failed to get worker 2: %v", err)
	}
	fmt.Printf("Got worker 2: %s\n", worker2)

	// Step 3: Do some work with the workers
	fmt.Println("\nDoing work with workers...")
	time.Sleep(2 * time.Second)

	// Step 4: Release individual worker (optional - they'll be released with session anyway)
	fmt.Printf("\nReleasing worker 1: %s\n", worker1)
	sessionPool.ReleaseWorker(worker1)

	// Step 5: Release the session (this will automatically release all remaining workers)
	fmt.Printf("\nReleasing session: %s (this will also release remaining workers)\n", sessionKey)
	sessionPool.Release(sessionKey)

	fmt.Println("\nDone!")
}
