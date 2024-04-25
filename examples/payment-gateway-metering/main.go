package main

import (
	"os"

	"github.com/streamingfast/dmetering"
	"github.com/streamingfast/payment-gateway/metering"
	"go.uber.org/zap"
)

func init() {
	metering.Register(
}

func main() {
	pluginDSN := os.ExpandEnv("paymentGateway://payment.gateway.streamingfast.io?network=eth-mainnet&token=${SF_API_TOKEN}")

	eventEmitter, err := dmetering.New(pluginDSN, zap.NewNop())
	if err != nil {
		panic(err)
	}

	_ = eventEmitter
}
