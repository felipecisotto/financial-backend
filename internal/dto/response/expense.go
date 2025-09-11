package response

import "time"

// ExpenseResponse representa os dados retornados de uma despesa
type ExpenseResponse struct {
	ID           string          `json:"id"`
	Description  string          `json:"description"`
	Amount       float64         `json:"amount"`
	Type         string          `json:"type"`
	BudgetID     *string         `json:"budget_id"`
	Budget       *BudgetResponse `json:"budget"`
	Recurrency   *string         `json:"recurrency"`
	Method       string          `json:"method"`
	Installments *int            `json:"installments"`
	DueDay       int             `json:"due_day"`
	StartDate    time.Time       `json:"start_date"`
	EndDate      *time.Time      `json:"end_date"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ListExpensesResponse representa a resposta da listagem de despesas
type ListExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int64             `json:"total"`
}
