package request

import "time"

// CreateIncomeRequest representa a requisição para criar uma receita
type CreateIncomeRequest struct {
	Description string     `json:"description" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	Type        string     `json:"type" binding:"required"`
	DueDay      int        `json:"due_day" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
}

// UpdateIncomeRequest representa a requisição para atualizar uma receita
type UpdateIncomeRequest struct {
	Description *string    `json:"description"`
	Amount      *float64   `json:"amount"`
	Type        *string    `json:"type"`
	Category    *string    `json:"category"`
	Date        *time.Time `json:"date"`
}

// ListIncomeParams representa os parâmetros para listar receitas
type ListIncomeParams struct {
	Type        string `form:"type"`
	Description string `form:"description"`
	Page        int64  `form:"page,default=1"`
	Limit       int64  `form:"limit,default=10"`
}
