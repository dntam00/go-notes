package monitor

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelMetric "go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"log"
	"net/http"
)

// docs: https://opentelemetry.io/docs/demo/

var (
	Meter           otelMetric.Meter
	DatabaseCommand otelMetric.Int64Histogram
)

func Start() {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Register prometheus exporter failed")
	}
	provider := sdkMetric.NewMeterProvider(sdkMetric.WithReader(exporter))
	Meter = provider.Meter("kaixin")

	DatabaseCommand, _ = Meter.Int64Histogram("database_command_latency", otelMetric.WithDescription("Command latency"))
}

func ServeMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

func RecordCommand(command string, latency int64) {
	DatabaseCommand.Record(context.Background(), latency, otelMetric.WithAttributes(
		attribute.Key("command").String(command)))
}
