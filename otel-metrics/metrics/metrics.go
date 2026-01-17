package metrics

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"

	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type config struct {
	endpoint        string
	applicationName string
	exportInterval  time.Duration
}

type Option func(*config)

func WithApplicationName(applicationName string) Option {
	return func(c *config) {
		c.applicationName = applicationName
	}
}

func WithExporterEndpoint(endpoint string) Option {
	return func(c *config) {
		c.endpoint = endpoint
	}
}

func WithExportInterval(interval time.Duration) Option {
	return func(c *config) {
		c.exportInterval = interval
	}
}

func NewProvider(opts ...Option) metric.MeterProvider {
	ctx := context.Background()

	c := new(config)
	for _, opt := range opts {
		opt(c)
	}

	// Step 1: Create OTLP gRPC exporter
	// This exporter PUSHES metrics to the OpenTelemetry Collector
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(c.endpoint),
		otlpmetricgrpc.WithInsecure(), // Use insecure for local dev
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(c.applicationName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Step 3: Create MeterProvider with periodic reader
	// The PeriodicReader controls HOW OFTEN metrics are exported
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter,
				sdkmetric.WithInterval(c.exportInterval),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider
}

func GetMeterProvider() metric.MeterProvider {
	return otel.GetMeterProvider()
}
