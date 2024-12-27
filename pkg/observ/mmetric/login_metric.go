package mmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

var loginFailCounter metric.Int64Counter

func init() {
	var err error
	loginFailCounter, err = meter.Int64Counter(
		"login.failed",
		metric.WithDescription("number of login failed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}
}

func AddLoginFailedCounter(ctx context.Context) {
	loginFailCounter.Add(ctx, 1, metric.WithAttributes(uniquePerNodeID))
}
