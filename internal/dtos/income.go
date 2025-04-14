package dtos

import "time"

// IncomeResponse representa a estrutura de resposta para receitas
type IncomeResponse struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	Type        string     `json:"type"`
	DueDay      int        `json:"due_day"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateIncomeRequest representa a requisição para criar uma receita
type CreateIncomeRequest struct {
	Description string     `json:"description" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	Type        string     `json:"type" binding:"required"`
	DueDay      int        `json:"due_day" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
}

type ListIncomeParams struct {
	Type        string `form:"type"`
	Description string `form:"description"`
	PageRequest
}

// UpdateIncomeRequest representa a requisição para atualizar uma receita
type UpdateIncomeRequest struct {
	Description *string    `json:"description"`
	Amount      *float64   `json:"amount"`
	Type        *string    `json:"type"`
	Category    *string    `json:"category"`
	Date        *time.Time `json:"date"`
}
