package observ

import (
	"context"
	"io"
	"os"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewTxtTracerProvider return a tracer provider with output
// traces.txt . Used for test in local
func NewTxtTracerProvider(mlog mlogger.Logger) (*trace.TracerProvider, func()) {
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		mlog.Error("fail create traces.txt", err)
	}
	tearsDown := make([]func(), 0, 2)
	tearsDown = append(tearsDown, func() {
		mlog.Info("close traces.txt file")
		if err := f.Close(); err != nil {
			mlog.Error("fail to close traces.txt file", err)
		}
	})

	exp, err := newExporter(f)
	if err != nil {
		mlog.Error("fail create exporter", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newTxtResource()),
	)
	tearsDown = append(tearsDown, func() {
		mlog.Info("shutdown tracer provider")
		if err := tp.Shutdown(context.Background()); err != nil {
			mlog.Error("fail to shutdown exporter", err)
		}
	})
	tearsDownfunc := func() {
		for _, ff := range tearsDown {
			ff()
		}
	}
	return tp, tearsDownfunc
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newTxtResource returns a resource describing this application.
func newTxtResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
