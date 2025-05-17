package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer             trace.Tracer
	meter              metric.Meter
	meterProvider      *sdkmetric.MeterProvider
	prometheusExporter *prometheus.Exporter
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
	var err2 error
	prometheusExporter, err2 = prometheus.New()
	if err2 != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err2)
	}

	// Create meter provider with the exporter
	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(prometheusExporter),
		sdkmetric.WithResource(res),
	)

	// Define o MeterProvider global
	// Isso garante que o middleware do Gin use o mesmo MeterProvider
	otel.SetMeterProvider(meterProvider)

	// Get meter from provider
	meter = meterProvider.Meter(serviceName)

	// Get tracer
	tracer = otel.Tracer(serviceName)

	// Configurar propagadores para rastreamento distribuído
	// Isso garante que o contexto de telemetria seja propagado entre serviços
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		// TraceContext propagador para rastreamento W3C
		propagation.TraceContext{},
		// Baggage propagador para metadados entre serviços
		propagation.Baggage{},
	))

	// Start metrics server
	go func() {
		mux := http.NewServeMux()
		// O handler do Prometheus expõe todas as métricas registradas
		// Incluindo aquelas criadas pelo middleware do Gin
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
	// Configuração para habilitar tanto traces quanto métricas
	// Com a versão 0.60.0+, o otelgin agora suporta métricas nativamente
	options := []otelgin.Option{
		// Usa explicitamente o MeterProvider para registrar métricas
		otelgin.WithMeterProvider(meterProvider),
	}

	// Este middleware agora irá gerar tanto traces quanto métricas
	// As métricas HTTP incluem:
	// - http.server.request.duration - tempo de processamento das requisições
	// - http.server.active_requests - número de requisições ativas
	// - http.server.request.body.size - tamanho do corpo da requisição
	// - http.server.response.body.size - tamanho do corpo da resposta
	return otelgin.Middleware("financial-backend", options...)
}

// Shutdown gracefully shuts down the telemetry system
func Shutdown(ctx context.Context) error {
	// Se necessário, podemos adicionar limpeza de recursos aqui no futuro
	return nil
}
