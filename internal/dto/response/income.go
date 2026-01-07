package response

import "time"

// IncomeResponse representa a estrutura de resposta para receitas
type IncomeResponse struct {
	ID                string     `json:"id"`
	Description       string     `json:"description"`
	Amount            float64    `json:"amount"`
	Type              string     `json:"type"`
	DueDay            int        `json:"due_day"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	DiscountMode      *string    `json:"discount_mode,omitempty"`
	HourlyRate        *float64   `json:"hourly_rate,omitempty"`
	HasSalaryTracking bool       `json:"has_salary_tracking"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
