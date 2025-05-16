package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	ginotel "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	tracer trace.Tracer
	meter  metric.Meter
)

// InitTelemetry initializes OpenTelemetry with Prometheus exporter
func InitTelemetry(serviceName string) error {
	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(meterProvider)

	// Get meter
	meter = meterProvider.Meter(serviceName)

	// Get tracer
	tracer = otel.Tracer(serviceName)

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9464", nil); err != nil {
			fmt.Printf("Error starting metrics server: %v\n", err)
		}
	}()

	return nil
}

// GetTracer returns the tracer instance
func GetTracer() trace.Tracer {
	return tracer
}

// GetMeter returns the meter instance
func GetMeter() metric.Meter {
	return meter
}

// GinMiddleware returns the OpenTelemetry middleware for Gin
func GinMiddleware() gin.HandlerFunc {
	return ginotel.Middleware("financial-backend")
}

// Shutdown gracefully shuts down the telemetry system
func Shutdown(ctx context.Context) error {
	return nil // No cleanup needed for Prometheus exporter
}
