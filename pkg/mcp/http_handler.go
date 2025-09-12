package mcp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPHandler provides an HTTP interface for MCP server
type HTTPHandler struct {
	server *Server
}

func NewHTTPHandler(server *Server) *HTTPHandler {
	return &HTTPHandler{
		server: server,
	}
}

// RegisterRoutes registers MCP HTTP endpoints with Gin router
func (h *HTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Basic HTTP endpoints for MCP
	router.GET("/info", h.serverInfo)
	router.POST("/call-tool", h.handleCallTool)
	router.GET("/list-tools", h.handleListTools)
	router.GET("/list-resources", h.handleListResources)
	router.GET("/list-prompts", h.handleListPrompts)
}

func (h *HTTPHandler) serverInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "financial-backend",
		"version": "v1.0.0",
		"type":    "mcp-server",
		"capabilities": gin.H{
			"tools":     true,
			"resources": true,
			"prompts":   true,
		},
		"protocol_version": "2024-11-05",
	})
}

// handleCallTool handles HTTP-based tool calls (non-streaming)
func (h *HTTPHandler) handleCallTool(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":              "HTTP tool calls not implemented",
		"message":            "Use streaming endpoint for tool execution",
		"streaming_endpoint": "/mcp/streaming/stream",
	})
}

// handleListTools returns available MCP tools
func (h *HTTPHandler) handleListTools(c *gin.Context) {
	// This would integrate with the registry to list available tools
	c.JSON(http.StatusOK, gin.H{
		"tools":   []gin.H{},
		"message": "Tools will be available after tool registration",
	})
}

// handleListResources returns available MCP resources
func (h *HTTPHandler) handleListResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"resources": []gin.H{},
		"message":   "Resources will be available after resource registration",
	})
}

// handleListPrompts returns available MCP prompts
func (h *HTTPHandler) handleListPrompts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"prompts": []gin.H{},
		"message": "Prompts will be available after prompt registration",
	})
}

// StartMCPInBackground starts the MCP server in a separate goroutine using STDIO transport
func (h *HTTPHandler) StartMCPInBackground(ctx context.Context) error {
	go func() {
		if err := h.server.Start(ctx); err != nil {
			// Log error but don't crash the main server
			// In a real implementation, you'd want proper error handling
		}
	}()
	return nil
}
