package auth

import "github.com/streamingfast/dmetrics"

var MetricSet = dmetrics.NewSet()
var ReissueCounter = MetricSet.NewCounter("auth_payment_gateway_reissue_counter", "How many times the 'reissue' method was called to get an updated JWT")
