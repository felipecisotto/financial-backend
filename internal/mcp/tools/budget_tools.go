package tools

import (
	"context"
	"fmt"
	"time"

	"financial-backend/internal/dto"
	"financial-backend/internal/dto/request"
	budgetUC "financial-backend/internal/usecases/budget"
	budgetMovementUC "financial-backend/internal/usecases/budget_movement"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BudgetTools provides MCP tools for budget management operations
type BudgetTools struct {
	budgetUseCase         budgetUC.UseCase
	budgetMovementUseCase budgetMovementUC.UseCase
}

// NewBudgetTools creates a new instance of BudgetTools
func NewBudgetTools(budgetUseCase budgetUC.UseCase, budgetMovementUseCase budgetMovementUC.UseCase) *BudgetTools {
	return &BudgetTools{
		budgetUseCase:         budgetUseCase,
		budgetMovementUseCase: budgetMovementUseCase,
	}
}

// CreateBudgetTool implements MCPTool interface for creating budgets
type CreateBudgetTool struct {
	budgetUseCase budgetUC.UseCase
}

func NewCreateBudgetTool(budgetUseCase budgetUC.UseCase) *CreateBudgetTool {
	return &CreateBudgetTool{
		budgetUseCase: budgetUseCase,
	}
}

func (t *CreateBudgetTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_budget",
		Description: "Create new budgets with validation",
		InputSchema: CreateBudgetSchema,
	}
}

func (t *CreateBudgetTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse optional end date
	var endDate *time.Time
	if endDateStr, exists := args["end_date"].(string); exists && endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %v", err)
		}
		endDate = &parsed
	}

	// Build request
	request := request.CreateBudgetRequest{
		Description: args["description"].(string),
		Amount:      args["amount"].(float64),
		EndDate:     endDate,
	}

	// Create budget
	response, err := t.budgetUseCase.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"budget":  response,
	}, nil
}

// GetBudgetsTool implements MCPTool interface for listing budgets
type GetBudgetsTool struct {
	budgetUseCase budgetUC.UseCase
}

func NewGetBudgetsTool(budgetUseCase budgetUC.UseCase) *GetBudgetsTool {
	return &GetBudgetsTool{
		budgetUseCase: budgetUseCase,
	}
}

func (t *GetBudgetsTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_budgets",
		Description: "Retrieve budgets with filtering options",
		InputSchema: GetBudgetsSchema,
	}
}

func (t *GetBudgetsTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := request.BudgetListParams{
		Page:  1,
		Limit: 10,
	}

	// Parse optional filters
	if description, exists := args["description"].(string); exists {
		params.Description = description
	}
	if status, exists := args["status"].(string); exists {
		params.Status = status
	}

	// Parse pagination
	if page, exists := args["page"].(float64); exists {
		params.Page = int64(page)
	}
	if limit, exists := args["limit"].(float64); exists {
		params.Limit = int64(limit)
	}

	// Get budgets
	result, err := t.budgetUseCase.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}

// GetBudgetByIdTool implements MCPTool interface for getting a specific budget
type GetBudgetByIdTool struct {
	budgetUseCase budgetUC.UseCase
}

func NewGetBudgetByIdTool(budgetUseCase budgetUC.UseCase) *GetBudgetByIdTool {
	return &GetBudgetByIdTool{
		budgetUseCase: budgetUseCase,
	}
}

func (t *GetBudgetByIdTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_budget_by_id",
		Description: "Get specific budget details by ID",
		InputSchema: GetBudgetByIdSchema,
	}
}

func (t *GetBudgetByIdTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	budget, err := t.budgetUseCase.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"budget":  budget,
	}, nil
}

// UpdateBudgetTool implements MCPTool interface for updating budgets
type UpdateBudgetTool struct {
	budgetUseCase budgetUC.UseCase
}

func NewUpdateBudgetTool(budgetUseCase budgetUC.UseCase) *UpdateBudgetTool {
	return &UpdateBudgetTool{
		budgetUseCase: budgetUseCase,
	}
}

func (t *UpdateBudgetTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "update_budget",
		Description: "Update existing budget information",
		InputSchema: UpdateBudgetSchema,
	}
}

func (t *UpdateBudgetTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	// Parse end date
	endDateStr, ok := args["end_date"].(string)
	if !ok {
		return nil, fmt.Errorf("end_date is required and must be a string")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %v", err)
	}

	// Build request
	request := &request.UpdateBudgetRequest{
		EndDate: endDate,
	}

	// Update budget
	response, err := t.budgetUseCase.Update(ctx, id, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"budget":  response,
	}, nil
}

// DeleteBudgetTool implements MCPTool interface for deleting budgets
type DeleteBudgetTool struct {
	budgetUseCase budgetUC.UseCase
}

func NewDeleteBudgetTool(budgetUseCase budgetUC.UseCase) *DeleteBudgetTool {
	return &DeleteBudgetTool{
		budgetUseCase: budgetUseCase,
	}
}

func (t *DeleteBudgetTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_budget",
		Description: "Delete budget records",
		InputSchema: DeleteBudgetSchema,
	}
}

func (t *DeleteBudgetTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := t.budgetUseCase.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Budget with ID %s deleted successfully", id),
	}, nil
}

// GetBudgetMovementsTool implements MCPTool interface for retrieving budget movements
type GetBudgetMovementsTool struct {
	budgetMovementUseCase budgetMovementUC.UseCase
}

func NewGetBudgetMovementsTool(budgetMovementUseCase budgetMovementUC.UseCase) *GetBudgetMovementsTool {
	return &GetBudgetMovementsTool{
		budgetMovementUseCase: budgetMovementUseCase,
	}
}

func (t *GetBudgetMovementsTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_budget_movements",
		Description: "Retrieve budget movement history",
		InputSchema: GetBudgetMovementsSchema,
	}
}

func (t *GetBudgetMovementsTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := dto.BudgetMovementParams{
		Page:  1,
		Limit: 10,
	}

	// Parse optional filters
	if budgetID, exists := args["budget_id"].(string); exists {
		params.BudgetId = budgetID
	}
	if movementType, exists := args["movement_type"].(string); exists {
		params.MovementType = movementType
	}
	if origin, exists := args["origin"].(string); exists {
		params.Origin = origin
	}
	if month, exists := args["month"].(float64); exists {
		params.Month = int(month)
	}
	if year, exists := args["year"].(float64); exists {
		params.Year = int(year)
	}

	// Parse pagination
	if page, exists := args["page"].(float64); exists {
		params.Page = int64(page)
	}
	if limit, exists := args["limit"].(float64); exists {
		params.Limit = int64(limit)
	}

	// Get budget movements
	result, err := t.budgetMovementUseCase.Find(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget movements: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}

// GetBudgetTools returns all budget-related MCP tools
func (bt *BudgetTools) GetBudgetTools() []interface{} {
	return []interface{}{
		NewCreateBudgetTool(bt.budgetUseCase),
		NewGetBudgetsTool(bt.budgetUseCase),
		NewGetBudgetByIdTool(bt.budgetUseCase),
		NewUpdateBudgetTool(bt.budgetUseCase),
		NewDeleteBudgetTool(bt.budgetUseCase),
		NewGetBudgetMovementsTool(bt.budgetMovementUseCase),
	}
}
