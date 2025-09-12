package mcp

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// MCPMiddleware provides middleware functions for MCP endpoints
type MCPMiddleware struct {
	tracer trace.Tracer
}

// NewMCPMiddleware creates a new MCP middleware instance
func NewMCPMiddleware() *MCPMiddleware {
	return &MCPMiddleware{
		tracer: otel.Tracer("mcp-middleware"),
	}
}

// LoggingMiddleware provides comprehensive logging for MCP endpoints
func (m *MCPMiddleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Get request info
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// Log the request
		log.Printf("MCP [%v] %s %s %d %v %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			log.Printf("MCP Errors: %v", c.Errors)
		}
	}
}

// ErrorHandlerMiddleware handles MCP-specific errors
func (m *MCPMiddleware) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("MCP Panic recovered: %v", r)

				// Create span for error tracking
				_, span := m.tracer.Start(c.Request.Context(), "mcp.error.panic_recovery")
				span.SetAttributes(
					attribute.String("error.type", "panic"),
					attribute.String("error.message", "panic recovered"),
				)
				span.End()

				c.JSON(500, gin.H{
					"error":   "Internal Server Error",
					"message": "An unexpected error occurred in MCP handler",
					"code":    "MCP_PANIC",
				})
			}
		}()

		c.Next()

		// Handle errors that were added to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Create span for error tracking
			_, span := m.tracer.Start(c.Request.Context(), "mcp.error.handler")
			span.SetAttributes(
				attribute.String("error.type", "gin_error"),
				attribute.String("error.message", err.Error()),
			)
			span.End()

			log.Printf("MCP Error: %v", err)

			// Return appropriate error response
			c.JSON(500, gin.H{
				"error":   "MCP Handler Error",
				"message": err.Error(),
				"code":    "MCP_ERROR",
			})
		}
	}
}

// TracingMiddleware adds OpenTelemetry tracing to MCP endpoints
func (m *MCPMiddleware) TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a new span
		ctx, span := m.tracer.Start(c.Request.Context(), "mcp.http.request")
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.route", c.FullPath()),
			attribute.String("user_agent", c.Request.UserAgent()),
		)

		// Update request context with span
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Set response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Mark span as error if needed
		if c.Writer.Status() >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	}
}

// CORSMiddleware handles CORS for MCP endpoints
func (m *MCPMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ConnectionLimitMiddleware limits concurrent streaming connections
func (m *MCPMiddleware) ConnectionLimitMiddleware(maxConnections int) gin.HandlerFunc {
	connectionCount := make(chan struct{}, maxConnections)

	return func(c *gin.Context) {
		select {
		case connectionCount <- struct{}{}:
			// Connection allowed
			defer func() { <-connectionCount }()
			c.Next()
		default:
			// Too many connections
			log.Printf("MCP connection limit reached: %d", maxConnections)
			c.JSON(429, gin.H{
				"error":   "Too Many Connections",
				"message": "Maximum number of concurrent MCP connections reached",
				"code":    "MCP_CONNECTION_LIMIT",
			})
			c.Abort()
		}
	}
}
