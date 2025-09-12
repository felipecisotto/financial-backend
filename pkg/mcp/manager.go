package mcp

import (
	"context"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

// Manager orchestrates all MCP components
type Manager struct {
	server           *Server
	streamingHandler *StreamingHandler
	httpHandler      *HTTPHandler
	middleware       *MCPMiddleware
	shutdownOnce     sync.Once
	isShutdown       bool
	mu               sync.RWMutex
}

// NewManager creates a new MCP manager with all components
func NewManager(config *Config) *Manager {
	// Create server
	server := NewServer(config)

	// Create handlers
	streamingHandler := NewStreamingHandler(server)
	httpHandler := NewHTTPHandler(server)
	middleware := NewMCPMiddleware()

	return &Manager{
		server:           server,
		streamingHandler: streamingHandler,
		httpHandler:      httpHandler,
		middleware:       middleware,
	}
}

// GetServer returns the MCP server instance
func (m *Manager) GetServer() *Server {
	return m.server
}

// GetStreamingHandler returns the streaming handler instance
func (m *Manager) GetStreamingHandler() *StreamingHandler {
	return m.streamingHandler
}

// GetHTTPHandler returns the HTTP handler instance
func (m *Manager) GetHTTPHandler() *HTTPHandler {
	return m.httpHandler
}

// RegisterRoutes registers all MCP routes with the Gin router
func (m *Manager) RegisterRoutes(router *gin.RouterGroup) {
	// Apply MCP middleware
	mcpGroup := router.Group("")
	mcpGroup.Use(m.middleware.TracingMiddleware())
	mcpGroup.Use(m.middleware.LoggingMiddleware())
	mcpGroup.Use(m.middleware.ErrorHandlerMiddleware())
	mcpGroup.Use(m.middleware.CORSMiddleware())
	mcpGroup.Use(m.middleware.ConnectionLimitMiddleware(100)) // Limit to 100 concurrent connections

	// Register streaming endpoints
	streamingGroup := mcpGroup.Group("/streaming")
	{
		m.streamingHandler.RegisterRoutes(streamingGroup)
	}

	// Register HTTP endpoints
	httpGroup := mcpGroup.Group("/http")
	{
		m.httpHandler.RegisterRoutes(httpGroup)
	}

	// Health check endpoint
	mcpGroup.GET("/health", m.healthCheck)

	log.Println("MCP routes registered successfully")
}

// healthCheck provides health status for MCP services
func (m *Manager) healthCheck(c *gin.Context) {
	m.mu.RLock()
	isShutdown := m.isShutdown
	m.mu.RUnlock()

	if isShutdown {
		c.JSON(503, gin.H{
			"status":  "shutdown",
			"message": "MCP service is shutting down",
		})
		return
	}

	activeConnections := m.streamingHandler.GetActiveConnections()

	c.JSON(200, gin.H{
		"status": "healthy",
		"mcp": gin.H{
			"server_ready":       true,
			"streaming_ready":    true,
			"active_connections": activeConnections,
			"server_name":        m.server.config.ServerName,
			"server_version":     m.server.config.ServerVersion,
		},
	})
}

// Shutdown gracefully shuts down all MCP components
func (m *Manager) Shutdown(ctx context.Context) error {
	var shutdownError error

	m.shutdownOnce.Do(func() {
		log.Println("Starting MCP manager shutdown...")

		// Mark as shutdown
		m.mu.Lock()
		m.isShutdown = true
		m.mu.Unlock()

		// Shutdown streaming handler
		if err := m.streamingHandler.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down streaming handler: %v", err)
			shutdownError = err
		}

		// Note: HTTP handler shutdown is handled by Gin server shutdown
		// No additional shutdown needed for HTTP handler

		log.Println("MCP manager shutdown complete")
	})

	return shutdownError
}

// IsShutdown returns whether the manager is in shutdown state
func (m *Manager) IsShutdown() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isShutdown
}

// GetStatus returns the current status of MCP components
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	isShutdown := m.isShutdown
	m.mu.RUnlock()

	return map[string]interface{}{
		"shutdown":           isShutdown,
		"active_connections": m.streamingHandler.GetActiveConnections(),
		"server_config": map[string]interface{}{
			"name":      m.server.config.ServerName,
			"version":   m.server.config.ServerVersion,
			"transport": m.server.config.Transport,
		},
	}
}
