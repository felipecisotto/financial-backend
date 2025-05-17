package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

// routeAttributeKey é a chave usada internamente para passar a rota pelo contexto
const routeAttributeKey = "otel-gin-route"

// GinMiddleware returns the OpenTelemetry middleware for Gin
func GinMiddleware() gin.HandlerFunc {
	// Criamos um middleware customizado que vai trabalhar em conjunto com o otelgin
	// para adicionar o atributo http.route às métricas
	ginHandler := func(c *gin.Context) {
		// Obter a rota completa do Gin
		route := c.FullPath()

		// Verificar se a rota existe
		if route != "" {
			// Armazenar a rota no contexto da requisição para uso posterior
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), routeAttributeKey, route))
		}

		// Continuar o processamento
		c.Next()
	}

	// Configuração do middleware do OpenTelemetry
	options := []otelgin.Option{
		// Usa explicitamente o MeterProvider para registrar métricas
		otelgin.WithMeterProvider(meterProvider),

		// Adiciona função para extrair o http.route das métricas
		otelgin.WithMetricAttributeFn(func(req *http.Request) []attribute.KeyValue {
			attrs := []attribute.KeyValue{}

			// Tenta extrair a rota do contexto da requisição
			if route, ok := req.Context().Value(routeAttributeKey).(string); ok && route != "" {
				attrs = append(attrs, attribute.String("http.route", route))
			}

			return attrs
		}),
	}

	// Cria o middleware do OpenTelemetry
	otelHandler := otelgin.Middleware("financial-backend", options...)

	// Retorna um handler que combina nosso middleware para capturar rotas
	// com o middleware do OpenTelemetry
	return func(c *gin.Context) {
		ginHandler(c)
		otelHandler(c)
	}
}

// Shutdown gracefully shuts down the telemetry system
func Shutdown(ctx context.Context) error {
	// Se necessário, podemos adicionar limpeza de recursos aqui no futuro
	return nil
}
