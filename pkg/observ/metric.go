package observ

import (
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

var meter = global.Meter("github.com/muchlist/moneymagnet")

var loginFailCounter, _ = meter.SyncInt64().Counter("login failed",
	instrument.WithDescription("number of login failed"),
	instrument.WithUnit(unit.Dimensionless),
)

func GetCounterLoginFailed() syncint64.Counter {
	return loginFailCounter
}
