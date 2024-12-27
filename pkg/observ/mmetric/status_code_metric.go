package mmetric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	endpointHitCounter metric.Int64Counter
	latencyHisto       metric.Int64Histogram
)

func init() {
	var err error
	endpointHitCounter, err = meter.Int64Counter(
		"http.request",
		metric.WithDescription("number of request per path"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}

	latencyHisto, err = meter.Int64Histogram(
		"response.latency",
		metric.WithDescription("latency of request"),
		metric.WithUnit("microseconds"),
	)
	if err != nil {
		panic(err)
	}
}

func AddEndpointHitCounter(ctx context.Context, code int, path string) {
	endpointHitCounter.Add(ctx, 1,
		metric.WithAttributes(
			uniquePerNodeID,
			attribute.String("method_path", path),
			attribute.Int("code", code),
		),
	)
}

func AddLatencyPerPath(ctx context.Context, durMicrosecond int64, path string) {
	latencyHisto.Record(ctx, durMicrosecond,
		metric.WithAttributes(
			uniquePerNodeID,
			attribute.String("method_path", path),
		),
	)
}
