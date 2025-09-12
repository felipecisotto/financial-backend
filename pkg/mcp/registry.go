package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry manages all MCP tools, resources, and prompts
type Registry struct {
	tools     []MCPTool
	resources []MCPResource
	prompts   []MCPPrompt
}

func NewRegistry() *Registry {
	return &Registry{
		tools:     make([]MCPTool, 0),
		resources: make([]MCPResource, 0),
		prompts:   make([]MCPPrompt, 0),
	}
}

func (r *Registry) RegisterTool(tool MCPTool) {
	r.tools = append(r.tools, tool)
	log.Printf("Registered MCP tool: %s", tool.GetTool().Name)
}

func (r *Registry) RegisterResource(resource MCPResource) {
	r.resources = append(r.resources, resource)
	log.Printf("Registered MCP resource: %s", resource.GetResource().Name)
}

func (r *Registry) RegisterPrompt(prompt MCPPrompt) {
	r.prompts = append(r.prompts, prompt)
	log.Printf("Registered MCP prompt: %s", prompt.GetPrompt().Name)
}

func (r *Registry) ConfigureServer(server *mcp.Server) {
	// Register all tools
	for _, tool := range r.tools {
		mcpTool := tool.GetTool()
		server.AddTool(mcpTool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Extract args from request
			args := make(map[string]interface{})
			// Note: Arguments handling will need proper JSON unmarshaling
			// For now, pass empty map - this will be improved in actual tool implementations

			result, err := tool.Execute(ctx, args)
			if err != nil {
				return nil, err
			}

			// Convert result to MCP response
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("%v", result),
					},
				},
			}, nil
		})
	}

	// Register all resources
	for _, resource := range r.resources {
		mcpResource := resource.GetResource()
		server.AddResource(mcpResource, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			content, err := resource.GetContent(ctx, req.Params.URI)
			if err != nil {
				return nil, err
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:      req.Params.URI,
						MIMEType: "text/plain",
						Text:     content,
					},
				},
			}, nil
		})
	}

	// Register all prompts
	for _, prompt := range r.prompts {
		mcpPrompt := prompt.GetPrompt()
		server.AddPrompt(mcpPrompt, func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			// Extract args from request
			args := make(map[string]interface{})
			if req.Params.Arguments != nil {
				for k, v := range req.Params.Arguments {
					args[k] = v
				}
			}

			messages, err := prompt.GetMessages(ctx, args)
			if err != nil {
				return nil, err
			}

			return &mcp.GetPromptResult{
				Messages: messages,
			}, nil
		})
	}

	log.Printf("Configured MCP server with %d tools, %d resources, %d prompts",
		len(r.tools), len(r.resources), len(r.prompts))
}
