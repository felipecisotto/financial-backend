package dtos

import "time"

type BudgetMovementRequest struct {
	BudgetId  string    `json:"budget_id"`
	Origin    string    `json:"origin"`
	Month     int       `json:"month"`
	Year      int       `json:"year"`
	Type      string    `json:"type"`
	Amount    int       `json:"amount"`
}

type BudgetMovementResponse struct {
	ID        string         `json:"id"`
	Budget    BudgetResponse `json:"budget"`
	Origin    string         `json:"origin"`
	Month     int            `json:"month"`
	Year      int            `json:"year"`
	Type      string         `json:"type"`
	Amount    int            `json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
}
