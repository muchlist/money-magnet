package mmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

// Login counter =========================
var loginFailCounter, _ = meter.SyncInt64().Counter("login failed",
	instrument.WithDescription("number of login failed"),
	instrument.WithUnit(unit.Dimensionless),
)

func GetCounterLoginFailed() syncint64.Counter {
	return loginFailCounter
}

func AddLoginFailedCounter(ctx context.Context) {
	// atrs := []attribute.KeyValue{
	// 	attribute.String("uid", uniqueDeploymentCode), // untuk membedakan antar node
	// }
	loginFailCounter.Add(ctx, 1, uniquePerNodeID)
}

// End of Login counter =====================
