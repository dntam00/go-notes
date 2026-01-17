// Package http provides an instrumented HTTP client with OpenTelemetry metrics.
//
// This package wraps the standard http.Client to automatically collect metrics:
// - Request count (by method, status code, host)
// - Request duration histogram
// - Active requests gauge
// - Request/Response size
//
// Usage:
//
//	client := otelhttp.NewClient(otelhttp.WithMeterProvider(provider))
//	resp, err := client.Get(ctx, "https://example.com")
package http

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	// Default meter name for HTTP client metrics
	defaultMeterName = "go-notes/otel-metrics/http"
)

// Client is an HTTP client instrumented with OpenTelemetry metrics.
type Client struct {
	client *http.Client

	// Metrics instruments
	requestCounter   metric.Int64Counter
	requestDuration  metric.Float64Histogram
	activeRequests   metric.Int64UpDownCounter
	requestBodySize  metric.Int64Counter
	responseBodySize metric.Int64Counter
}

// Option configures the Client.
type Option func(*clientConfig)

type clientConfig struct {
	meterProvider metric.MeterProvider
	meterName     string
	httpClient    *http.Client
}

// WithMeterProvider sets the meter provider to use for metrics.
// If not set, the global meter provider is used.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return func(c *clientConfig) {
		c.meterProvider = provider
	}
}

// WithMeterName sets the meter name for metrics instrumentation.
func WithMeterName(name string) Option {
	return func(c *clientConfig) {
		c.meterName = name
	}
}

// WithHTTPClient sets the underlying HTTP client to use.
// If not set, http.DefaultClient is used.
func WithHTTPClient(client *http.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

// NewClient creates a new instrumented HTTP client.
func NewClient(opts ...Option) (*Client, error) {
	cfg := &clientConfig{
		meterName:  defaultMeterName,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Use global meter provider if not specified
	meterProvider := cfg.meterProvider
	if meterProvider == nil {
		meterProvider = otel.GetMeterProvider()
	}

	meter := meterProvider.Meter(cfg.meterName)

	// Create metrics instruments
	requestCounter, err := meter.Int64Counter(
		"http.client.request.count",
		metric.WithDescription("Total number of HTTP requests made"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"http.client.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"http.client.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	requestBodySize, err := meter.Int64Counter(
		"http.client.request.body.size",
		metric.WithDescription("Total bytes sent in HTTP request bodies"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	responseBodySize, err := meter.Int64Counter(
		"http.client.response.body.size",
		metric.WithDescription("Total bytes received in HTTP response bodies"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:           cfg.httpClient,
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
		activeRequests:   activeRequests,
		requestBodySize:  requestBodySize,
		responseBodySize: responseBodySize,
	}, nil
}

// Do sends an HTTP request and returns an HTTP response with metrics collection.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Common attributes for all metrics
	attrs := []attribute.KeyValue{
		attribute.String("http.method", req.Method),
		attribute.String("http.host", req.URL.Host),
		attribute.String("http.scheme", req.URL.Scheme),
	}

	// Track active requests
	c.activeRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
	defer c.activeRequests.Add(ctx, -1, metric.WithAttributes(attrs...))

	// Track request body size if available
	if req.ContentLength > 0 {
		c.requestBodySize.Add(ctx, req.ContentLength, metric.WithAttributes(attrs...))
	}

	// Execute request
	resp, err := c.client.Do(req.WithContext(ctx))

	// Calculate duration
	duration := time.Since(start).Seconds()

	// Add status code to attributes
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}
	attrs = append(attrs, attribute.Int("http.status_code", statusCode))

	// Add error attribute if request failed
	if err != nil {
		attrs = append(attrs, attribute.Bool("http.error", true))
	}

	// Record metrics
	c.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	c.requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))

	// Track response body size
	if resp != nil && resp.ContentLength > 0 {
		c.responseBodySize.Add(ctx, resp.ContentLength, metric.WithAttributes(attrs...))
	}

	return resp, err
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(ctx, req)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(ctx, req)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

// Head performs a HEAD request.
func (c *Client) Head(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

// StatusClass returns the status class (1xx, 2xx, 3xx, 4xx, 5xx) for a status code.
func StatusClass(statusCode int) string {
	return strconv.Itoa(statusCode/100) + "xx"
}
