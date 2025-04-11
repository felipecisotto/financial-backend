package dtos

import (
	"time"
)

// CreateExpenseRequest representa os dados necessários para criar uma despesa
type CreateExpenseRequest struct {
	Description  string     `json:"description" binding:"required"`
	Amount       float64    `json:"amount" binding:"required"`
	Type         string     `json:"type" binding:"required"`
	BudgetID     *string    `json:"budget_id"`
	DueDate 	 *time.Time  `json:"due_date"`
	Installments *int       `json:"installments"`
}

// ExpenseResponse representa os dados retornados de uma despesa
type ExpenseResponse struct {
	ID           string     `json:"id"`
	Description  string     `json:"description"`
	Amount       float64    `json:"amount"`
	Type         string     `json:"type"`
	BudgetID     *string    `json:"budget_id"`
	StartDate    *time.Time `json:"start_date"`
	DueDate      *time.Time `json:"end_date"`
}

// ListExpensesResponse representa a resposta da listagem de despesas
type ListExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int64             `json:"total"`
}

// ListExpensesByMonthRequest representa os parâmetros para listar despesas por mês
type ListExpensesByMonthRequest struct {
	Month time.Time `json:"month" binding:"required"`
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
