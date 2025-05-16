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
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	ginotel "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	tracer trace.Tracer
	meter  metric.Meter
)

// InitTelemetry initializes OpenTelemetry with Prometheus exporter
func InitTelemetry(serviceName string) error {
	// Create a resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}
	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}


	// Create meter provider with the exporter
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Get meter from provider

	meter = meterProvider.Meter(serviceName)

	// Get tracer
	tracer = otel.Tracer(serviceName)

	// Start metrics server
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:    ":9464",
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
