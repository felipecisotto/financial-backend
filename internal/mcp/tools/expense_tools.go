package tools

import (
	"context"
	"fmt"
	"time"

	"financial-backend/internal/dto"
	"financial-backend/internal/usecases/expense"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExpenseTools provides MCP tools for expense management operations
type ExpenseTools struct {
	expenseUseCase expense.UseCase
}

// NewExpenseTools creates a new instance of ExpenseTools
func NewExpenseTools(expenseUseCase expense.UseCase) *ExpenseTools {
	return &ExpenseTools{
		expenseUseCase: expenseUseCase,
	}
}

// CreateExpenseTool implements MCPTool interface for creating expenses
type CreateExpenseTool struct {
	expenseUseCase expense.UseCase
}

func NewCreateExpenseTool(expenseUseCase expense.UseCase) *CreateExpenseTool {
	return &CreateExpenseTool{
		expenseUseCase: expenseUseCase,
	}
}

func (t *CreateExpenseTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_expense",
		Description: "Create a new expense with validation",
		InputSchema: CreateExpenseSchema,
	}
}

func (t *CreateExpenseTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	// Parse optional budget_id
	var budgetID *string
	if budgetIDStr, exists := args["budget_id"].(string); exists && budgetIDStr != "" {
		budgetID = &budgetIDStr
	}

	// Parse optional recurrency
	var recurrency *string
	if recurrencyStr, exists := args["recurrency"].(string); exists && recurrencyStr != "" {
		recurrency = &recurrencyStr
	}

	// Parse optional installments
	var installments *int
	if installmentsFloat, exists := args["installments"].(float64); exists {
		installmentsInt := int(installmentsFloat)
		installments = &installmentsInt
	}

	// Build request
	request := &dto.CreateExpenseRequest{
		Description:  args["description"].(string),
		Amount:       args["amount"].(float64),
		Type:         args["type"].(string),
		BudgetID:     budgetID,
		Recurrency:   recurrency,
		Method:       args["method"].(string),
		Installments: installments,
		DueDay:       int(args["due_day"].(float64)),
		StartDate:    startDate,
		EndDate:      endDate,
	}

	// Create expense
	response, err := t.expenseUseCase.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"expense": response,
	}, nil
}

// GetExpensesTool implements MCPTool interface for listing expenses
type GetExpensesTool struct {
	expenseUseCase expense.UseCase
}

func NewGetExpensesTool(expenseUseCase expense.UseCase) *GetExpensesTool {
	return &GetExpensesTool{
		expenseUseCase: expenseUseCase,
	}
}

func (t *GetExpensesTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_expenses",
		Description: "Retrieve expenses with filtering options",
		InputSchema: GetExpensesSchema,
	}
}

func (t *GetExpensesTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	request := &dto.ListExpensesRequest{
		Page:  1,
		Limit: 10,
	}

	// Parse optional filters
	if description, exists := args["description"].(string); exists {
		request.Description = description
	}
	if expenseType, exists := args["type"].(string); exists {
		request.Type = expenseType
	}
	if category, exists := args["category"].(string); exists {
		request.Category = category
	}
	if budgetID, exists := args["budget_id"].(string); exists {
		request.BudgetID = budgetID
	}
	if recurrency, exists := args["recurrency"].(string); exists {
		request.Recurrency = recurrency
	}
	if method, exists := args["method"].(string); exists {
		request.Method = method
	}

	// Parse pagination
	if page, exists := args["page"].(float64); exists {
		request.Page = int64(page)
	}
	if limit, exists := args["limit"].(float64); exists {
		request.Limit = int64(limit)
	}

	// Get expenses
	result, err := t.expenseUseCase.List(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}

// GetExpenseByIdTool implements MCPTool interface for getting a specific expense
type GetExpenseByIdTool struct {
	expenseUseCase expense.UseCase
}

func NewGetExpenseByIdTool(expenseUseCase expense.UseCase) *GetExpenseByIdTool {
	return &GetExpenseByIdTool{
		expenseUseCase: expenseUseCase,
	}
}

func (t *GetExpenseByIdTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_expense_by_id",
		Description: "Get specific expense details by ID",
		InputSchema: GetExpenseByIdSchema,
	}
}

func (t *GetExpenseByIdTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	expense, err := t.expenseUseCase.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"expense": expense,
	}, nil
}

// DeleteExpenseTool implements MCPTool interface for deleting expenses
type DeleteExpenseTool struct {
	expenseUseCase expense.UseCase
}

func NewDeleteExpenseTool(expenseUseCase expense.UseCase) *DeleteExpenseTool {
	return &DeleteExpenseTool{
		expenseUseCase: expenseUseCase,
	}
}

func (t *DeleteExpenseTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_expense",
		Description: "Delete expense records",
		InputSchema: DeleteExpenseSchema,
	}
}

func (t *DeleteExpenseTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}

	err := t.expenseUseCase.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete expense: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Expense with ID %s deleted successfully", id),
	}, nil
}

