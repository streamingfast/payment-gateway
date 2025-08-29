package main

import (
	"fmt"
	"os"

	"github.com/streamingfast/dmetering"
	"github.com/streamingfast/payment-gateway/metering"
	"go.uber.org/zap"
)

func init() {
	// Register tgm:// as a valid metering plugin, refers to Register documentation for extra details
	metering.Register()
}

func main() {
	checkEnv("API_KEY")

	// The ${API_TOKEN} is expanded automatically by the metering plugin via your defind environment variables
	// (as well as any other variables)!
	//
	// The [network] query parameter is required and should match the well-know network identifier.
	pluginDSN := "tgm://abp.thegraph.market?network=eth-mainnet&api_key=${API_KEY}"

	eventEmitter, err := dmetering.New(pluginDSN, zap.NewNop())
	if err != nil {
		panic(err)
	}

	_ = eventEmitter
}

func checkEnv(key string) {
	if _, exists := os.LookupEnv(key); !exists {
		panic(fmt.Errorf("missing required environment variable %q", key))
	}
}
