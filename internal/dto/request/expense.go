package request

import "time"

// CreateExpenseRequest representa os dados necessários para criar uma despesa
type CreateExpenseRequest struct {
	Description  string     `json:"description" binding:"required"`
	Amount       float64    `json:"amount" binding:"required"`
	Type         string     `json:"type" binding:"required"`
	BudgetID     *string    `json:"budget_id"`
	Recurrency   *string    `json:"recurrency"`
	Method       string     `json:"method" binding:"required"`
	Installments *int       `json:"installments"`
	DueDay       int        `json:"due_day" binding:"required"`
	StartDate    time.Time  `json:"start_date" binding:"required"`
	EndDate      *time.Time `json:"end_date"`
}

// UpdateExpenseRequest representa os dados necessários para atualizar uma despesa
type UpdateExpenseRequest struct {
	Description  *string    `json:"description"`
	Amount       *float64   `json:"amount"`
	Category     *string    `json:"category"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Installments *int       `json:"installments"`
	DueDate      *time.Time `json:"due_date"`
	StatementDay *int       `json:"statement_day"`
}

// ListExpensesRequest representa os parâmetros para listar despesas
type ListExpensesRequest struct {
	Description string `form:"description"`
	Type        string `form:"type"`
	Category    string `form:"category"`
	BudgetID    string `form:"budget_id"`
	Recurrency  string `form:"recurrency"`
	Method      string `form:"method"`
	Page        int64  `form:"page,default=1"`
	Limit       int64  `form:"limit,default=10"`
}