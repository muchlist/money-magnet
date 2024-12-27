package observ

import (
	"context"
	"crypto/tls"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/credentials"
)

type exporterFactory struct {
	ctx context.Context
	opt Option
}

func newExporterFactory(ctx context.Context, opt Option) *exporterFactory {
	return &exporterFactory{
		ctx: ctx,
		opt: opt,
	}
}

func (f *exporterFactory) createTraceExporter() (*otlptrace.Exporter, error) {
	secureOption := f.getSecureOption()
	return otlptrace.New(
		f.ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(f.opt.CollectorURL),
			otlptracegrpc.WithHeaders(f.opt.Headers),
		),
	)
}

func (f *exporterFactory) createMetricExporter() (metric.Exporter, error) {
	secureOption := f.getMetricSecureOption()
	return otlpmetricgrpc.New(f.ctx,
		otlpmetricgrpc.WithEndpoint(f.opt.CollectorURL),
		secureOption,
	)
}

func (f *exporterFactory) getSecureOption() otlptracegrpc.Option {
	if f.opt.Insecure {
		return otlptracegrpc.WithInsecure()
	}
	// Gunakan sistem sertifikat default atau specify sertifikat
	creds := credentials.NewTLS(&tls.Config{})
	return otlptracegrpc.WithTLSCredentials(creds)
}

func (f *exporterFactory) getMetricSecureOption() otlpmetricgrpc.Option {
	if f.opt.Insecure {
		return otlpmetricgrpc.WithInsecure()
	}
	return otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
}
