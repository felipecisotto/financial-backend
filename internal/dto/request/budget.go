package request

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

// BudgetListParams representa os parâmetros para listar orçamentos
type BudgetListParams struct {
	Description string `form:"description"`
	Status      string `form:"status"`
	Page        int64  `form:"page,default=1"`
	Limit       int64  `form:"limit,default=10"`
}