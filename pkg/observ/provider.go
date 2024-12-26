package observ

import (
	"context"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(ctx context.Context, opt Option, log mlogger.Logger) func(context.Context) error {
	factory := newExporterFactory(ctx, opt)
	exporter, err := factory.createTraceExporter()
	if err != nil {
		log.Error("failed create exporter", err)
		panic(err)
	}

	resources, err := createResource(ctx, opt.ServiceName)
	if err != nil {
		log.Error("could not set resources", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	))

	return exporter.Shutdown
}

func InitMeter(ctx context.Context, opt Option, log mlogger.Logger) func(context.Context) error {
	factory := newExporterFactory(ctx, opt)
	exporter, err := factory.createMetricExporter()
	if err != nil {
		log.Error("failed create otlpmetricgrpc exporter", err)
	}

	resources, err := createResource(ctx, opt.ServiceName)
	if err != nil {
		log.Error("failed create resource", err)
	}

	reader := metric.NewPeriodicReader(exporter)

	opts := []metric.Option{
		metric.WithReader(reader),
		metric.WithResource(resources),
	}

	if customBucketsView := createCustomBucketsView(); customBucketsView != nil {
		opts = append(opts, metric.WithView(customBucketsView))
	}

	provider := metric.NewMeterProvider(opts...)

	otel.SetMeterProvider(provider)
	return provider.Shutdown
}

func createResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.node", "303030"),
			attribute.String("library.language", "go"),
		),
	)
}

func createCustomBucketsView() metric.View {
	return metric.NewView(
		metric.Instrument{
			Name: "response.*",
			Kind: metric.InstrumentKindHistogram,
		},
		metric.Stream{
			Aggregation: metric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{25_000, 50_000, 100_000, 250_000, 500_000, 1_000_000, 2_500_000, 5_000_000},
			},
		},
	)
}
