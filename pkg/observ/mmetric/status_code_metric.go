package mmetric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
)

// endpoint hit counter =========================
var endpointHitCounter, _ = meter.SyncInt64().Counter("endpoint hit",
	instrument.WithDescription("number of endpoint hit per path"),
	instrument.WithUnit(unit.Dimensionless),
)

func AddEndpointHitCounter(ctx context.Context, code int, path string) {
	atrs := []attribute.KeyValue{
		attribute.String("uid", uniqueDeploymentCode),
		attribute.String("method_path", path),
		attribute.Int("code", code),
	}
	endpointHitCounter.Add(ctx, 1, atrs...)
}

// End of endpoint hit counter =====================

// status code counter =========================
var statusCodeCounter, _ = meter.SyncInt64().Counter("response code",
	instrument.WithDescription("number of response code"),
	instrument.WithUnit(unit.Dimensionless),
)

func AddStatusCodeCounter(ctx context.Context, code int) {
	atrs := []attribute.KeyValue{
		attribute.String("uid", uniqueDeploymentCode),
		attribute.Int("code", code),
	}
	statusCodeCounter.Add(ctx, 1, atrs...)
}

// End of status code counter =====================
