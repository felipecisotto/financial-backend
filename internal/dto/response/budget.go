package response

import "time"

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
