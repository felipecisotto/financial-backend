package request

// BudgetMovementRequest representa a requisição para criar um movimento de orçamento
type BudgetMovementRequest struct {
	BudgetId string `json:"budget_id"`
	Origin   string `json:"origin"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Type     string `json:"type"`
	Amount   int    `json:"amount"`
}

// BudgetMovementParams representa os parâmetros para listar movimentos
type BudgetMovementParams struct {
	BudgetId     string `form:"budget_id"`
	MovementType string `form:"movement_type"`
	Origin       string `form:"origin"`
	Month        int    `form:"month"`
	Year         int    `form:"year"`
	Page         int64  `form:"page,default=1"`
	Limit        int64  `form:"limit,default=10"`
}