package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcpServer *mcp.Server
}

func NewServer() *Server {
	implementation := &mcp.Implementation{
		Name:    "financial-backend-mcp",
		Version: "1.0.0",
	}

	mcpServer := mcp.NewServer(implementation, nil)

	return &Server{
		mcpServer: mcpServer,
	}
}
func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
