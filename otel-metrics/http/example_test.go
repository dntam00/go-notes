package http_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	otelhttp "go-notes/otel-metrics/http"
)

func Example() {
	ctx := context.Background()

	// Create OTLP exporter
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create resource
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("http-client-example"),
		),
	)

	// Create MeterProvider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter,
				sdkmetric.WithInterval(10*time.Second),
			),
		),
	)
	defer meterProvider.Shutdown(ctx)

	// Set as global
	otel.SetMeterProvider(meterProvider)

	// Create instrumented HTTP client
	client, err := otelhttp.NewClient(
		otelhttp.WithMeterProvider(meterProvider),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Make requests - metrics are automatically collected
	resp, err := client.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		log.Printf("request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %d\n", resp.StatusCode)
	// Output: Status: 200
}
