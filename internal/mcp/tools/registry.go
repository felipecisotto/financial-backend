package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	budgetUC "financial-backend/internal/usecases/budget"
	budgetMovementUC "financial-backend/internal/usecases/budget_movement"
	dashboardUC "financial-backend/internal/usecases/dashboard"
	expenseUC "financial-backend/internal/usecases/expense"
	incomeUC "financial-backend/internal/usecases/income"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolRegistry manages registration and discovery of all MCP tools
type ToolRegistry struct {
	mu               sync.RWMutex
	tools            map[string]MCPTool
	categories       map[ToolCategory][]string
	metadata         map[string]MCPToolMetadata
	initialized      bool
	expenseUseCase   expenseUC.UseCase
	incomeUseCase    incomeUC.UseCase
	budgetUseCase    budgetUC.UseCase
	budgetMovementUC budgetMovementUC.UseCase
	dashboardUseCase dashboardUC.UseCase
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:      make(map[string]MCPTool),
		categories: make(map[ToolCategory][]string),
		metadata:   make(map[string]MCPToolMetadata),
	}
}

// SetUseCases sets the use case dependencies for tool initialization
func (r *ToolRegistry) SetUseCases(
	expenseUC expenseUC.UseCase,
	incomeUC incomeUC.UseCase,
	budgetUC budgetUC.UseCase,
	budgetMovementUC budgetMovementUC.UseCase,
	dashboardUC dashboardUC.UseCase,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.expenseUseCase = expenseUC
	r.incomeUseCase = incomeUC
	r.budgetUseCase = budgetUC
	r.budgetMovementUC = budgetMovementUC
	r.dashboardUseCase = dashboardUC
}

// Initialize registers all available tools with their dependencies
func (r *ToolRegistry) Initialize() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		return fmt.Errorf("registry already initialized")
	}

	// Register expense tools
	if err := r.registerExpenseTools(); err != nil {
		return fmt.Errorf("failed to register expense tools: %v", err)
	}

	// Register income tools
	if err := r.registerIncomeTools(); err != nil {
		return fmt.Errorf("failed to register income tools: %v", err)
	}

	// Register budget tools
	if err := r.registerBudgetTools(); err != nil {
		return fmt.Errorf("failed to register budget tools: %v", err)
	}

	// Register dashboard tools
	if err := r.registerDashboardTools(); err != nil {
		return fmt.Errorf("failed to register dashboard tools: %v", err)
	}

	r.initialized = true
	log.Printf("Tool registry initialized with %d tools", len(r.tools))
	return nil
}

// ConfigureMCPServer configures the MCP server with all registered tools
func (r *ToolRegistry) ConfigureMCPServer(server *mcp.Server) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.initialized {
		return fmt.Errorf("registry not initialized")
	}

	for name, tool := range r.tools {
		mcpTool := tool.GetTool()
		server.AddTool(mcpTool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Printf("Tool handler called for: %s", name)

			// Extract args from request
			args := make(map[string]interface{})
			if req.Params.Arguments != nil {
				// Arguments is json.RawMessage, unmarshal it
				if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
					log.Printf("Warning: Failed to unmarshal arguments: %v", err)
				}
			}

			result, err := tool.Execute(ctx, args)
			if err != nil {
				return nil, err
			}

			// Serialize result to proper JSON format
			serializedResult := serializeToolResult(result)

			// Convert result to MCP response
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: serializedResult,
					},
				},
			}, nil
		})

		log.Printf("Configured MCP server with tool: %s", name)
	}

	log.Printf("MCP server configured with %d tools", len(r.tools))
	return nil
}

// registerToolUnsafe registers a tool without acquiring lock (internal use)
func (r *ToolRegistry) registerToolUnsafe(tool MCPTool, metadata MCPToolMetadata) error {
	if _, exists := r.tools[metadata.Name]; exists {
		return fmt.Errorf("tool %s already registered", metadata.Name)
	}

	r.tools[metadata.Name] = tool
	r.metadata[metadata.Name] = metadata

	// Add to category
	if r.categories[metadata.Category] == nil {
		r.categories[metadata.Category] = make([]string, 0)
	}
	r.categories[metadata.Category] = append(r.categories[metadata.Category], metadata.Name)

	log.Printf("Registered MCP tool: %s (category: %s)", metadata.Name, metadata.Category)
	return nil
}

// registerExpenseTools registers all expense-related tools
func (r *ToolRegistry) registerExpenseTools() error {
	if r.expenseUseCase == nil {
		return fmt.Errorf("expense use case not provided")
	}

	expenseTools := []struct {
		tool     MCPTool
		metadata MCPToolMetadata
	}{
		{
			tool: NewCreateExpenseTool(r.expenseUseCase),
			metadata: MCPToolMetadata{
				Name:        "create_expense",
				Description: "Create new expenses with validation",
				Category:    CategoryExpense,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetExpensesTool(r.expenseUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_expenses",
				Description: "Retrieve expenses with filtering options",
				Category:    CategoryExpense,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetExpenseByIdTool(r.expenseUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_expense_by_id",
				Description: "Get specific expense details by ID",
				Category:    CategoryExpense,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewDeleteExpenseTool(r.expenseUseCase),
			metadata: MCPToolMetadata{
				Name:        "delete_expense",
				Description: "Delete expense records",
				Category:    CategoryExpense,
				Version:     "1.0.0",
			},
		},
	}

	for _, toolInfo := range expenseTools {
		if err := r.registerToolUnsafe(toolInfo.tool, toolInfo.metadata); err != nil {
			return err
		}
	}

	return nil
}

// registerIncomeTools registers all income-related tools
func (r *ToolRegistry) registerIncomeTools() error {
	if r.incomeUseCase == nil {
		return fmt.Errorf("income use case not provided")
	}

	incomeTools := []struct {
		tool     MCPTool
		metadata MCPToolMetadata
	}{
		{
			tool: NewCreateIncomeTool(r.incomeUseCase),
			metadata: MCPToolMetadata{
				Name:        "create_income",
				Description: "Create new income records with validation",
				Category:    CategoryIncome,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetIncomesTool(r.incomeUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_incomes",
				Description: "Retrieve income records with filtering options",
				Category:    CategoryIncome,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetIncomeByIdTool(r.incomeUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_income_by_id",
				Description: "Get specific income details by ID",
				Category:    CategoryIncome,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewUpdateIncomeTool(r.incomeUseCase),
			metadata: MCPToolMetadata{
				Name:        "update_income",
				Description: "Update existing income information",
				Category:    CategoryIncome,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewDeleteIncomeTool(r.incomeUseCase),
			metadata: MCPToolMetadata{
				Name:        "delete_income",
				Description: "Delete income records",
				Category:    CategoryIncome,
				Version:     "1.0.0",
			},
		},
	}

	for _, toolInfo := range incomeTools {
		if err := r.registerToolUnsafe(toolInfo.tool, toolInfo.metadata); err != nil {
			return err
		}
	}

	return nil
}

// registerBudgetTools registers all budget-related tools
func (r *ToolRegistry) registerBudgetTools() error {
	if r.budgetUseCase == nil || r.budgetMovementUC == nil {
		return fmt.Errorf("budget or budget movement use case not provided")
	}

	budgetTools := []struct {
		tool     MCPTool
		metadata MCPToolMetadata
	}{
		{
			tool: NewCreateBudgetTool(r.budgetUseCase),
			metadata: MCPToolMetadata{
				Name:        "create_budget",
				Description: "Create new budgets with validation",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetBudgetsTool(r.budgetUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_budgets",
				Description: "Retrieve budgets with filtering options",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetBudgetByIdTool(r.budgetUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_budget_by_id",
				Description: "Get specific budget details by ID",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewUpdateBudgetTool(r.budgetUseCase),
			metadata: MCPToolMetadata{
				Name:        "update_budget",
				Description: "Update existing budget information",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewDeleteBudgetTool(r.budgetUseCase),
			metadata: MCPToolMetadata{
				Name:        "delete_budget",
				Description: "Delete budget records",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetBudgetMovementsTool(r.budgetMovementUC),
			metadata: MCPToolMetadata{
				Name:        "get_budget_movements",
				Description: "Retrieve budget movement history",
				Category:    CategoryBudget,
				Version:     "1.0.0",
			},
		},
	}

	for _, toolInfo := range budgetTools {
		if err := r.registerToolUnsafe(toolInfo.tool, toolInfo.metadata); err != nil {
			return err
		}
	}

	return nil
}

// registerDashboardTools registers all dashboard-related tools
func (r *ToolRegistry) registerDashboardTools() error {
	if r.dashboardUseCase == nil {
		return fmt.Errorf("dashboard use case not provided")
	}

	dashboardTools := []struct {
		tool     MCPTool
		metadata MCPToolMetadata
	}{
		{
			tool: NewGetMonthlySummaryTool(r.dashboardUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_monthly_summary",
				Description: "Get monthly financial breakdown",
				Category:    CategoryDashboard,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetBudgetUtilizationTool(r.dashboardUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_budget_utilization",
				Description: "Get budget usage and remaining amounts",
				Category:    CategoryDashboard,
				Version:     "1.0.0",
			},
		},
		{
			tool: NewGetIncomeVsExpensesTool(r.dashboardUseCase),
			metadata: MCPToolMetadata{
				Name:        "get_income_vs_expenses",
				Description: "Compare income against expenses",
				Category:    CategoryDashboard,
				Version:     "1.0.0",
			},
		},
	}

	for _, toolInfo := range dashboardTools {
		if err := r.registerToolUnsafe(toolInfo.tool, toolInfo.metadata); err != nil {
			return err
		}
	}

	return nil
}

// serializeToolResult handles automatic JSON serialization for all tool results
func serializeToolResult(result interface{}) string {
	// If result is already a string, assume it's pre-serialized JSON
	if str, ok := result.(string); ok {
		return str
	}

	// Otherwise, serialize to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Warning: Failed to serialize tool result: %v", err)
		// Fallback to null JSON
		return "null"
	}

	return string(jsonData)
}

// Global registry instance
var globalRegistry *ToolRegistry
var registryOnce sync.Once

// GetGlobalRegistry returns the global tool registry instance
func GetGlobalRegistry() *ToolRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewToolRegistry()
	})
	return globalRegistry
}
