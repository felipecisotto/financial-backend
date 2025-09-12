package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPTool represents a tool that can be registered with MCP
type MCPTool interface {
	GetTool() *mcp.Tool
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// ToolCategory represents the category of a tool
type ToolCategory string

const (
	CategoryExpense   ToolCategory = "expense"
	CategoryIncome    ToolCategory = "income"
	CategoryBudget    ToolCategory = "budget"
	CategoryDashboard ToolCategory = "dashboard"
)

// MCPToolMetadata contains metadata about a tool
type MCPToolMetadata struct {
	Name        string
	Description string
	Category    ToolCategory
	Version     string
}
