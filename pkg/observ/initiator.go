package observ

import (
	"context"
	"time"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

type Option struct {
	ServiceName string
	// Without https:// (example: localhost:4317)
	CollectorURL string
	ApiKey       string
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
	return exporter.Shutdown
}

func getTraceExporter(ctx context.Context, opt Option) (*otlptrace.Exporter, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if opt.Insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}
	return otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(opt.CollectorURL),
			otlptracegrpc.WithHeaders(
				map[string]string{
					"api-key": opt.ApiKey,
				},
			),
		),
	)
}

func getResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
}

// InitMeter...
func InitMeter(ctx context.Context, opt Option, log mlogger.Logger) func(context.Context) error {

	exporter, err := getMeterExporter(ctx, opt.Insecure, opt.CollectorURL)
	if err != nil {
		log.Error("failed to create metric exporter", err)
	}

	resources, err := getResource(ctx, opt.ServiceName)
	if err != nil {
		log.Error("failed create resource", err)
	}

	cont := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries([]float64{1, 2, 5, 10, 20, 50}),
			),
			temporalitySelector,
		),
		controller.WithResource(resources),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(2*time.Second),
	)

	err = cont.Start(ctx)
	if err != nil {
		log.Error("failed to start controller", err)
	}

	global.SetMeterProvider(cont)

	return exporter.Shutdown
}

type newRelicTemporalitySelector struct{}

func (s newRelicTemporalitySelector) TemporalityFor(desc *sdkapi.Descriptor, kind aggregation.Kind) aggregation.Temporality {
	if desc.InstrumentKind() == sdkapi.CounterInstrumentKind ||
		// The Go SDK doesn't support Async Observers with Delta temporality yet.
		// To avoid errors, use cumulative for Async Counters, which NR will interpret as gauges.
		// desc.InstrumentKind() == sdkapi.CounterObserverInstrumentKind ||
		desc.InstrumentKind() == sdkapi.HistogramInstrumentKind {
		return aggregation.DeltaTemporality
	}
	return aggregation.CumulativeTemporality
}

var temporalitySelector = newRelicTemporalitySelector{}

func getMeterExporter(ctx context.Context, insecure bool, address string) (*otlpmetric.Exporter, error) {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlpmetricgrpc.WithInsecure()
	}
	return otlpmetric.New(
		ctx,
		otlpmetricgrpc.NewClient(
			secureOption,
			otlpmetricgrpc.WithEndpoint(address),
		),
		otlpmetric.WithMetricAggregationTemporalitySelector(temporalitySelector),
	)
}
