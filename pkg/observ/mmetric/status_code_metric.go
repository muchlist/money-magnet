package mmetric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
)

// http request counter =========================
var endpointHitCounter, _ = meter.SyncInt64().Counter("http.request",
	instrument.WithDescription("number of request per path"),
	instrument.WithUnit(unit.Dimensionless),
)

func AddEndpointHitCounter(ctx context.Context, code int, path string) {
	atrs := []attribute.KeyValue{
		uniquePerNodeID,
		attribute.String("method_path", path),
		attribute.Int("code", code),
	}
	endpointHitCounter.Add(ctx, 1, atrs...)
}

// End of endpoint hit counter =====================

// Latency histogram =========================
var latencyHisto, _ = meter.SyncInt64().Histogram("response.latency",
	instrument.WithDescription("latency of request"),
	instrument.WithUnit("microseconds"),
)

func AddLatencyPerPath(ctx context.Context, durMicrosecond int64, path string) {
	atrs := []attribute.KeyValue{
		uniquePerNodeID,
		attribute.String("method_path", path),
	}
	latencyHisto.Record(ctx, durMicrosecond, atrs...)
}

// End of Latency histogram =====================
