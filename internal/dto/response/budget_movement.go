package response

import "time"

// BudgetMovementResponse representa a resposta com os dados de um movimento
type BudgetMovementResponse struct {
	ID                string         `json:"id"`
	Budget            BudgetResponse `json:"budget"`
	Origin            string         `json:"origin"`
	OriginDescription *string        `json:"origin_description"`
	Month             int            `json:"month"`
	Year              int            `json:"year"`
	Type              string         `json:"type"`
	Amount            int            `json:"amount"`
	CreatedAt         time.Time      `json:"created_at"`
}