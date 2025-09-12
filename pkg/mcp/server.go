package mcp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcpServer *mcp.Server
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

	return &Server{
		mcpServer: mcpServer,
		config:    config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	log.Println("Starting MCP server...")
	transport := &mcp.StreamableClientTransport{}
	return s.mcpServer.Run(ctx, transport)
}

func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
