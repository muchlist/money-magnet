package observ

import (
	"context"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Option struct {
	ServiceName string
	// Without https:// (example: localhost:4317)
	CollectorURL string
	Headers      map[string]string
	Insecure     bool
}

// InitTracer ...
func InitTracer(ctx context.Context, opt Option, log mlogger.Logger) func(context.Context) error {
	exporter, err := getTraceExporter(ctx, opt)
	if err != nil {
		log.Error("failed create exporter", err)
		panic(err)
	}
	resources, err := getResource(ctx, opt.ServiceName)
	if err != nil {
		log.Error("could not set resources", err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithBatcher(exporter),
			trace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

	return exporter.Shutdown
}

// InitMeter...
func InitMeter(ctx context.Context, opt Option, log mlogger.Logger) func(context.Context) error {

	exporter, err := getOtlpMetricExporter(ctx, opt)
	if err != nil {
		log.Error("failed create otlpmetricgrpc exporter", err)
	}

	resources, err := getResource(ctx, opt.ServiceName)
	if err != nil {
		log.Error("failed create resource", err)
	}

	provider := metric.NewMeterProvider(
		// metric.WithReader(exporter),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resources),
	)

	global.SetMeterProvider(provider)

	return provider.Shutdown
}

func getResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.node", "303030"),
			attribute.String("library.language", "go"),
		),
	)
}
