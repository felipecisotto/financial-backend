package dtos

import (
	"time"
)

// ExpenseDTO representa os dados necessários para criar uma despesa
type ExpenseDTO struct {
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

// ExpenseResponse representa os dados retornados de uma despesa
type ExpenseResponse struct {
	ID string `json:"id"`
	ExpenseDTO
}

// ListExpensesResponse representa a resposta da listagem de despesas
type ListExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int64             `json:"total"`
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
	Type     *string    `json:"type"`
	Category *string    `json:"category"`
	BudgetID *string    `json:"budget_id"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
}
