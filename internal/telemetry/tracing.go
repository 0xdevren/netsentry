package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerOptions configures the OpenTelemetry tracer provider.
type TracerOptions struct {
	// ServiceName is the OTel service name attribute.
	ServiceName string
	// ServiceVersion is the OTel service version attribute.
	ServiceVersion string
	// Enabled controls whether tracing is active.
	Enabled bool
}

// InitTracer initialises the OpenTelemetry tracer with a stdout exporter.
// It returns a shutdown function that flushes and closes the provider.
func InitTracer(opts TracerOptions) (trace.Tracer, func(context.Context) error, error) {
	if !opts.Enabled {
		noop := otel.Tracer(opts.ServiceName)
		return noop, func(_ context.Context) error { return nil }, nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opts.ServiceName),
			semconv.ServiceVersion(opts.ServiceVersion),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(provider)

	tracer := provider.Tracer(opts.ServiceName)
	return tracer, provider.Shutdown, nil
}

// Tracer returns the global OTel tracer for the given instrumentation scope.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
