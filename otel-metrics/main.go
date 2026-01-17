// OpenTelemetry Metrics - OTLP Push Model Example
//
// This example demonstrates how OpenTelemetry metrics work with OTLP push:
// - Metrics are periodically PUSHED to an OpenTelemetry Collector
// - The SDK handles batching and export intervals automatically
// - No HTTP server needed (unlike Prometheus pull model)
//
// Prerequisites:
//   - OpenTelemetry Collector running at localhost:4317 (gRPC)
//
// Run with: go run main.go

package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	otelhttp "go-notes/otel-metrics/http"
	"go-notes/otel-metrics/metrics"
)

const (
	exportInterval    = 1 * time.Second
	collectorEndpoint = "127.0.0.1:4317"
)

func main() {
	_ = metrics.NewProvider(
		metrics.WithApplicationName("kaixin-metrics-example"),
		metrics.WithExporterEndpoint(collectorEndpoint),
		metrics.WithExportInterval(exportInterval),
	)

	go testHttp()
	go testRandom()
	select {}
}

func testHttp() {
	for {
		provider := metrics.GetMeterProvider()
		client, err := otelhttp.NewClient(
			otelhttp.WithMeterProvider(provider),
		)
		if err != nil {
			log.Fatal(err)
		}

		provider.Meter("kaixin-http")

		resp, err := client.Get(context.Background(), "https://httpbin.org/get")
		if err != nil {
			log.Printf("request failed: %v", err)
			return
		}
		_ = resp.Body.Close()
		time.Sleep(2 * time.Second)
	}
}

func testRandom() {
	ctx := context.Background()
	meter := metrics.GetMeterProvider().Meter("kaixin-random")

	activeConnections, _ := meter.Int64UpDownCounter("active_connections")

	for {
		delta := rand.Intn(3) - 1
		if delta != 0 {
			activeConnections.Add(ctx, int64(delta))
		}
		time.Sleep(2 * time.Second)
	}
}
