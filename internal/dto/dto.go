package dto

// Re-export commonly used types for easier imports
import (
	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
)

// Request types
type (
	CreateExpenseRequest  = request.CreateExpenseRequest
	UpdateExpenseRequest  = request.UpdateExpenseRequest
	ListExpensesRequest   = request.ListExpensesRequest
	CreateBudgetRequest   = request.CreateBudgetRequest
	UpdateBudgetRequest   = request.UpdateBudgetRequest
	BudgetListParams      = request.BudgetListParams
	CreateIncomeRequest   = request.CreateIncomeRequest
	UpdateIncomeRequest   = request.UpdateIncomeRequest
	ListIncomeParams      = request.ListIncomeParams
	BudgetMovementRequest = request.BudgetMovementRequest
	BudgetMovementParams  = request.BudgetMovementParams
	SummaryQueryParams    = request.SummaryQueryParams
)

// Response types
type (
	ExpenseResponse          = response.ExpenseResponse
	ListExpensesResponse     = response.ListExpensesResponse
	BudgetResponse           = response.BudgetResponse
	IncomeResponse           = response.IncomeResponse
	BudgetMovementResponse   = response.BudgetMovementResponse
	SummaryView              = response.SummaryView
	SummaryBudgetUtilization = response.SummaryBudgetUtilization
)

// Legacy types for backward compatibility (remove after refactoring)
type (
	ExpenseDTO = CreateExpenseRequest // ExpenseDTO maps to CreateExpenseRequest
)

// Common types
type PageRequest struct {
	Page  int64 `form:"page,default=1"`
	Limit int64 `form:"limit,default=10"`
}
