package observ

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/credentials"
)

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
			otlptracegrpc.WithHeaders(opt.Headers),
		),
	)
}

func getOtlpMetricExporter(ctx context.Context, opt Option) (metric.Exporter, error) {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if opt.Insecure {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(opt.CollectorURL),
		secureOption,
	)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func getPrometheuspMetricExporter(ctx context.Context, opt Option) (*prometheus.Exporter, error) {
	// need to import "go.opentelemetry.io/otel/exporters/prometheus"
	// need to enable expose /metrics using prometheus http middlware

	return prometheus.New()
}
