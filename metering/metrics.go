package metering

import "github.com/streamingfast/dmetrics"

var MetricSet = dmetrics.NewSet()
var DroppedEventCounter = MetricSet.NewCounter("metering_payment_gateway_dropped_event_counter", "Counter of drop metering events from the metering paymentGateway:// service")
var MeteringGRPCErrCounter = MetricSet.NewCounter("metering_payment_gateway_err_counter", "Counter of gRPC errors received reporting to the Payment Gateway service")
var MeteringGRPCRetryCounter = MetricSet.NewCounter("metering_payment_gateway_retry_counter", "Counter of gRPC retries and trying to report to the Payment Gateway service")
