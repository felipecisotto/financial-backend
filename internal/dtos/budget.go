package dtos

import "time"

// CreateBudgetRequest representa a requisição para criar um orçamento
type CreateBudgetRequest struct {
	Description string     `json:"description" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
}

// UpdateBudgetRequest representa a requisição para atualizar um orçamento
type UpdateBudgetRequest struct {
	EndDate time.Time `json:"end_date" binding:"required"`
}

// BudgetResponse representa a resposta com os dados de um orçamento
type BudgetResponse struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BudgetSummaryResponse representa o resumo do orçamento
type BudgetSummaryResponse struct {
	TotalBudgeted   float64 `json:"total_budgeted"`
	TotalSpent      float64 `json:"total_spent"`
	Remaining       float64 `json:"remaining"`
	PercentageSpent float64 `json:"percentage_spent"`
}
