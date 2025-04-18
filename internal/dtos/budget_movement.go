package dtos

import "time"

type BudgetMovementRequest struct {
	BudgetId string `json:"budget_id"`
	Origin   string `json:"origin"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Type     string `json:"type"`
	Amount   int    `json:"amount"`
}
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

type BudgetMovementParams struct {
	BudgetId     string `form:"budget_id"`
	MovementType string `form:"movement_type"`
	Origin       string `form:"origin"`
	Month        int    `form:"month"`
	Year         int    `form:"year"`
	PageRequest
}
