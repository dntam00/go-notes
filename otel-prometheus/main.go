// OpenTelemetry Metrics - Prometheus Exporter (Pull Model) Example
//
// This example demonstrates how OpenTelemetry metrics work with Prometheus:
// - Metrics are exposed via HTTP endpoint (/metrics)
// - Prometheus PULLS (scrapes) metrics from this endpoint
// - The scrape interval is controlled by Prometheus, not the application
//
// Key difference from OTLP push:
// - OTLP Push: Application pushes metrics to collector on a schedule
// - Prometheus Pull: Prometheus scrapes metrics from application on its schedule
//
// Run with: go run main.go
// Then visit: http://localhost:2112/metrics

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	// HTTP server port for Prometheus to scrape
	metricsPort = ":8000"
)

func main() {
	ctx := context.Background()

	// Step 1: Create Prometheus exporter
	// This exporter exposes metrics in Prometheus format at /metrics
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("failed to create prometheus exporter: %v", err)
	}

	// Step 2: Create resource (identifies this application)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-prometheus-example"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Step 3: Create MeterProvider with Prometheus exporter
	// Note: No PeriodicReader needed - Prometheus controls scrape timing
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter), // Prometheus exporter IS the reader
	)
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("error shutting down meter provider: %v", err)
		}
	}()

	// Step 4: Set as global meter provider
	otel.SetMeterProvider(meterProvider)

	// Step 5: Create a meter (namespace for metrics)
	meter := otel.Meter("example.com/metrics")

	// Step 6: Create a Counter metric
	requestCounter, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	// Step 7: Create a Gauge metric (via UpDownCounter)
	activeConnections, err := meter.Int64UpDownCounter(
		"active_connections",
		metric.WithDescription("Number of active connections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		log.Fatalf("failed to create gauge: %v", err)
	}

	// Step 8: Start HTTP server for Prometheus to scrape
	// The promhttp.Handler() serves metrics in Prometheus text format
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		fmt.Printf("Metrics endpoint: http://localhost%s/metrics\n", metricsPort)
		if err := http.ListenAndServe(metricsPort, nil); err != nil {
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()

	fmt.Printf("OpenTelemetry + Prometheus Example Started\n")
	fmt.Printf("==========================================\n")
	fmt.Printf("Metrics endpoint: http://localhost%s/metrics\n", metricsPort)
	fmt.Printf("\nPrometheus PULLS metrics when it scrapes (you control scrape_interval in prometheus.yml)\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// Simulate some work and record metrics
	for {
		// Simulate incoming requests (increment counter)
		requestCounter.Add(ctx, 1, metric.WithAttributes())
		fmt.Printf("[%s] Recorded 1 request\n", time.Now().Format("15:04:05"))

		// Simulate connection changes (random up/down)
		delta := rand.Intn(3) - 1 // -1, 0, or 1
		if delta != 0 {
			activeConnections.Add(ctx, int64(delta))
			fmt.Printf("[%s] Active connections changed by %d\n", time.Now().Format("15:04:05"), delta)
		}

		time.Sleep(2 * time.Second)
	}
}
