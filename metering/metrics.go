package metering

import "github.com/streamingfast/dmetrics"

var MetricSet = dmetrics.NewSet()
var DroppedEventCounter = MetricSet.NewCounter("metering_payment_gateway_dropped_event_counter", "Counter of drop metering events from the metering paymentGateway:// service")
var MeteringGRPCErrCounter = MetricSet.NewCounter("metering_payment_gateway_err_counter", "Counter of gRPC errors received from the Payment Gateway service")
