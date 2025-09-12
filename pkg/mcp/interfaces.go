package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPTool represents a tool that can be registered with the MCP server
type MCPTool interface {
	// GetTool returns the MCP tool definition
	GetTool() *mcp.Tool

	// Execute handles the tool execution
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// MCPResource represents a resource that can be accessed via MCP
type MCPResource interface {
	// GetResource returns the MCP resource definition
	GetResource() *mcp.Resource

	// GetContent returns the resource content
	GetContent(ctx context.Context, uri string) (string, error)
}

// MCPPrompt represents a prompt template available via MCP
type MCPPrompt interface {
	// GetPrompt returns the MCP prompt definition
	GetPrompt() *mcp.Prompt

	// GetMessages returns the prompt messages
	GetMessages(ctx context.Context, args map[string]interface{}) ([]*mcp.PromptMessage, error)
}
