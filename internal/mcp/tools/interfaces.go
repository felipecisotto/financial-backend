package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPTool represents a tool that can be registered with the MCP server
// This interface should be implemented by all MCP tools
type MCPTool interface {
	// GetTool returns the MCP tool definition
	GetTool() *mcp.Tool

	// Execute handles the tool execution
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// MCPToolMetadata contains metadata about a tool
type MCPToolMetadata struct {
	Name         string
	Description  string
	Category     ToolCategory
	Version      string
	Dependencies []string
}

// ToolCategory represents different categories of tools
type ToolCategory string

const (
	CategoryExpense   ToolCategory = "expense"
	CategoryIncome    ToolCategory = "income"
	CategoryBudget    ToolCategory = "budget"
	CategoryDashboard ToolCategory = "dashboard"
	CategoryGeneral   ToolCategory = "general"
)

// MCPToolWithMetadata combines a tool with its metadata
type MCPToolWithMetadata interface {
	MCPTool
	GetMetadata() MCPToolMetadata
}

// ToolInitializer is a function that creates and initializes a tool with its dependencies
type ToolInitializer func(deps *ToolDependencies) (MCPTool, error)

// ToolDependencies contains all dependencies that tools might need
type ToolDependencies struct {
	// Add dependency interfaces here as needed
	// Example: ExpenseUseCase, IncomeUseCase, etc.
}
