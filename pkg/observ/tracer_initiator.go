package observ

import (
	"context"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

type Option struct {
	ServiceName  string
	CollectorURL string
	Insecure     bool
}

func InitTracer(opt Option, log mlogger.Logger) func(context.Context) error {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !opt.Insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(opt.CollectorURL),
		),
	)

	if err != nil {
		log.Error("failed create exporter", err)
		panic(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", opt.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
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
