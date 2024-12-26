package observ

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var moneyMagnetTracer = otel.Tracer("github.com/muchlist/moneymagnet")

func GetTracer() trace.Tracer {
	return moneyMagnetTracer
}
