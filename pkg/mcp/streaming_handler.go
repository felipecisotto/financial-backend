package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// StreamingHandler handles MCP streaming connections using SSE
type StreamingHandler struct {
	server         *Server
	connections    map[string]*ConnectionContext
	sessionManager *SessionManager
	mu             sync.RWMutex
	tracer         trace.Tracer
}

// ConnectionContext manages individual streaming connections
type ConnectionContext struct {
	ID        string
	StartTime time.Time
	Context   context.Context
	Cancel    context.CancelFunc
}

// NewStreamingHandler creates a new MCP streaming handler
func NewStreamingHandler(server *Server) *StreamingHandler {
	return &StreamingHandler{
		server:         server,
		connections:    make(map[string]*ConnectionContext),
		sessionManager: NewSessionManager(),
		tracer:         otel.Tracer("mcp-streaming-handler"),
	}
}

// HandleWebSocket handles WebSocket connections for MCP (future implementation)
func (sh *StreamingHandler) HandleWebSocket(c *gin.Context) {
	_, span := sh.tracer.Start(c.Request.Context(), "mcp.streaming.websocket.handle")
	defer span.End()

	// TODO: Implement WebSocket support when needed
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "WebSocket not implemented yet",
		"message": "Use SSE endpoint for streaming",
	})
}

// HandleStreamingInfo provides information about active streaming connections
func (sh *StreamingHandler) HandleStreamingInfo(c *gin.Context) {
	_, span := sh.tracer.Start(c.Request.Context(), "mcp.streaming.info")
	defer span.End()

	sh.mu.RLock()
	defer sh.mu.RUnlock()

	connections := make([]gin.H, 0, len(sh.connections))
	for sessionID, conn := range sh.connections {
		connections = append(connections, gin.H{
			"connection_id": conn.ID,
			"session_id":    sessionID,
			"start_time":    conn.StartTime,
			"duration":      time.Since(conn.StartTime).String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"active_connections": len(sh.connections),
		"connections":        connections,
		"server_info": gin.H{
			"name":    "financial-backend-mcp",
			"version": "v1.0.0",
		},
	})
}

// Shutdown gracefully closes all streaming connections
func (sh *StreamingHandler) Shutdown(ctx context.Context) error {
	_, span := sh.tracer.Start(ctx, "mcp.streaming.shutdown")
	defer span.End()

	log.Println("Shutting down MCP streaming handler...")

	sh.mu.Lock()
	defer sh.mu.Unlock()

	// Cancel all active connections
	for sessionID, conn := range sh.connections {
		log.Printf("Closing MCP streaming connection %s (session: %s)", conn.ID, sessionID)
		conn.Cancel()
	}

	// Wait for connections to close or timeout
	timeout := time.NewTimer(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer timeout.Stop()
	defer ticker.Stop()

	for len(sh.connections) > 0 {
		select {
		case <-timeout.C:
			log.Printf("Timeout waiting for %d connections to close", len(sh.connections))
			return fmt.Errorf("timeout waiting for connections to close")
		case <-ticker.C:
			// Continue checking
		}
	}

	log.Println("MCP streaming handler shutdown complete")
	return nil
}

// GetActiveConnections returns the number of active streaming connections
func (sh *StreamingHandler) GetActiveConnections() int {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return len(sh.connections)
}

// generateConnectionID creates a unique connection identifier
func (sh *StreamingHandler) generateConnectionID() string {
	return fmt.Sprintf("mcp-conn-%d", time.Now().UnixNano())
}

// RegisterRoutes registers streaming endpoints with the Gin router
func (sh *StreamingHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Create the MCP SSE handler once
	sseHandler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		connectionID := sh.generateConnectionID()
		log.Printf("Creating MCP server for connection %s", connectionID)
		return sh.server.GetMCPServer()
	})

	// Wrap the MCP handler with session management middleware
	sessionManagedHandler := sh.sessionManager.SessionMiddleware(sseHandler)

	// SSE endpoint for streaming MCP connections
	router.GET("/stream", gin.WrapH(sessionManagedHandler))
	router.POST("/stream", gin.WrapH(sessionManagedHandler))

	// The MCP SDK might create additional endpoints internally, so we need to handle them
	router.POST("/stream/*path", gin.WrapH(sessionManagedHandler))

	// WebSocket endpoint (future implementation)
	router.GET("/ws", sh.HandleWebSocket)

	// Info endpoint for connection status
	router.GET("/info", sh.HandleStreamingInfo)
}
