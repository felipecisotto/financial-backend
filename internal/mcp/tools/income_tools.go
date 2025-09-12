package tools

import (
	"context"
	"fmt"
	"time"

	"financial-backend/internal/dto"
	"financial-backend/internal/usecases/income"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// IncomeTools provides MCP tools for income management operations
type IncomeTools struct {
	incomeUseCase income.UseCase
}

// NewIncomeTools creates a new instance of IncomeTools
func NewIncomeTools(incomeUseCase income.UseCase) *IncomeTools {
	return &IncomeTools{
		incomeUseCase: incomeUseCase,
	}
}

// CreateIncomeTool implements MCPTool interface for creating income records
type CreateIncomeTool struct {
	incomeUseCase income.UseCase
}

func NewCreateIncomeTool(incomeUseCase income.UseCase) *CreateIncomeTool {
	return &CreateIncomeTool{
		incomeUseCase: incomeUseCase,
	}
}

func (t *CreateIncomeTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_income",
		Description: "Create new income records with validation",
		InputSchema: CreateIncomeSchema,
	}
}

func (t *CreateIncomeTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse start date
	startDateStr, ok := args["start_date"].(string)
	if !ok {
		return nil, fmt.Errorf("start_date is required and must be a string")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %v", err)
	}

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
	request := &dto.CreateIncomeRequest{
		Description: args["description"].(string),
		Amount:      args["amount"].(float64),
		Type:        args["type"].(string),
		DueDay:      int(args["due_day"].(float64)),
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Create income
	response, err := t.incomeUseCase.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create income: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"income":  response,
	}, nil
}

// GetIncomesTool implements MCPTool interface for listing income records
type GetIncomesTool struct {
	incomeUseCase income.UseCase
}

func NewGetIncomesTool(incomeUseCase income.UseCase) *GetIncomesTool {
	return &GetIncomesTool{
		incomeUseCase: incomeUseCase,
	}
}

func (t *GetIncomesTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_incomes",
		Description: "Retrieve income records with filtering options",
		InputSchema: GetIncomesSchema,
	}
}

func (t *GetIncomesTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	params := dto.ListIncomeParams{
		Page:  1,
		Limit: 10,
	}

	// Parse optional filters
	if description, exists := args["description"].(string); exists {
		params.Description = description
	}
	if incomeType, exists := args["type"].(string); exists {
		params.Type = incomeType
	}

	// Parse pagination
	if page, exists := args["page"].(float64); exists {
		params.Page = int64(page)
	}
	if limit, exists := args["limit"].(float64); exists {
		params.Limit = int64(limit)
	}

	// Get incomes
	result, err := t.incomeUseCase.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get incomes: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}

// GetIncomeByIdTool implements MCPTool interface for getting a specific income record
type GetIncomeByIdTool struct {
	incomeUseCase income.UseCase
}

func NewGetIncomeByIdTool(incomeUseCase income.UseCase) *GetIncomeByIdTool {
	return &GetIncomeByIdTool{
		incomeUseCase: incomeUseCase,
	}
}

func (t *GetIncomeByIdTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_income_by_id",
		Description: "Get specific income details by ID",
		InputSchema: GetIncomeByIdSchema,
	}
}

func (t *GetIncomeByIdTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	income, err := t.incomeUseCase.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"income":  income,
	}, nil
}

// UpdateIncomeTool implements MCPTool interface for updating income records
type UpdateIncomeTool struct {
	incomeUseCase income.UseCase
}

func NewUpdateIncomeTool(incomeUseCase income.UseCase) *UpdateIncomeTool {
	return &UpdateIncomeTool{
		incomeUseCase: incomeUseCase,
	}
}

func (t *UpdateIncomeTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "update_income",
		Description: "Update existing income information",
		InputSchema: UpdateIncomeSchema,
	}
}

func (t *UpdateIncomeTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	request := &dto.UpdateIncomeRequest{}

	// Parse optional fields
	if description, exists := args["description"].(string); exists {
		request.Description = &description
	}
	if amount, exists := args["amount"].(float64); exists {
		request.Amount = &amount
	}
	if incomeType, exists := args["type"].(string); exists {
		request.Type = &incomeType
	}
	if category, exists := args["category"].(string); exists {
		request.Category = &category
	}
	if dateStr, exists := args["date"].(string); exists {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		request.Date = &date
	}

	// Update income
	response, err := t.incomeUseCase.Update(ctx, id, request)
	if err != nil {
		return nil, fmt.Errorf("failed to update income: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"income":  response,
	}, nil
}

// DeleteIncomeTool implements MCPTool interface for deleting income records
type DeleteIncomeTool struct {
	incomeUseCase income.UseCase
}

func NewDeleteIncomeTool(incomeUseCase income.UseCase) *DeleteIncomeTool {
	return &DeleteIncomeTool{
		incomeUseCase: incomeUseCase,
	}
}

func (t *DeleteIncomeTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_income",
		Description: "Delete income records",
		InputSchema: DeleteIncomeSchema,
	}
}

func (t *DeleteIncomeTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := t.incomeUseCase.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete income: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Income with ID %s deleted successfully", id),
	}, nil
}

// GetIncomeTools returns all income-related MCP tools
func (it *IncomeTools) GetIncomeTools() []interface{} {
	return []interface{}{
		NewCreateIncomeTool(it.incomeUseCase),
		NewGetIncomesTool(it.incomeUseCase),
		NewGetIncomeByIdTool(it.incomeUseCase),
		NewUpdateIncomeTool(it.incomeUseCase),
		NewDeleteIncomeTool(it.incomeUseCase),
	}
}
