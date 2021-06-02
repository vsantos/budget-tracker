package observability

import (
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	zipkinExporter "go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// Config o
type Config struct {
	TracerProviders ProvidersConfig
}

// ProvidersConfig will
type ProvidersConfig struct {
	ServiceName string
	ZipkinURL   string
	JaegerURL   string
}

// Providers will
type Providers struct {
	Stdout *tracesdk.TracerProvider
	Zipkin *tracesdk.TracerProvider
	Jaeger *tracesdk.TracerProvider
}

// InitTracerProviders will return a struct with both providers: jaeger and stdout
func InitTracerProviders(c ProvidersConfig) (p Providers, err error) {
	resourceAttributes := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(c.ServiceName),
	)

	zp, err := ZipkinTracerProvider(resourceAttributes, c.ZipkinURL)
	if err != nil {
		return p, err
	}

	jp, err := JaegerTracerProvider(resourceAttributes, c.JaegerURL)
	if err != nil {
		return p, err
	}

	sp, err := StdoutTracerProvider(resourceAttributes)
	if err != nil {
		return p, err
	}

	p = Providers{
		Zipkin: zp,
		Jaeger: jp,
		Stdout: sp,
	}

	return p, nil
}

// JaegerTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func JaegerTracerProvider(resourceAttributes *resource.Resource, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter

	exp, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resourceAttributes),
	)

	return tp, nil
}

// ZipkinTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Zipkin exporter
func ZipkinTracerProvider(resourceAttributes *resource.Resource, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := zipkinExporter.NewRawExporter(url)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resourceAttributes),
	)

	return tp, nil
}

// StdoutTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the stdout exporter
func StdoutTracerProvider(resourceAttributes *resource.Resource) (*tracesdk.TracerProvider, error) {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(resourceAttributes),
	)

	return tp, nil
}
