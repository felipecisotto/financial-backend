package mcp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcpServer *mcp.Server
	registry  *Registry
	config    *Config
}

func NewServer(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	implementation := &mcp.Implementation{
		Name:    config.ServerName,
		Version: config.ServerVersion,
	}

	mcpServer := mcp.NewServer(implementation, nil)
	registry := NewRegistry()

	return &Server{
		mcpServer: mcpServer,
		registry:  registry,
		config:    config,
	}
}

func (s *Server) GetRegistry() *Registry {
	return s.registry
}

func (s *Server) RegisterHandlers() {
	s.registry.ConfigureServer(s.mcpServer)
}

func (s *Server) Start(ctx context.Context) error {
	s.RegisterHandlers()

	log.Println("Starting MCP server...")

	// Use STDIO transport for now - can be configured later
	transport := &mcp.StdioTransport{}
	return s.mcpServer.Run(ctx, transport)
}

func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
