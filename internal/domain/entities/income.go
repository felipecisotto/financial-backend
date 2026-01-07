package entities

import (
	"time"
)

// Income representa a tabela de receitas
type Income struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Description  string     `json:"description"`
	Amount       float64    `json:"amount"`
	Type         string     `json:"type"`
	StartDate    time.Time  `json:"start_date"`
	DueDay       int        `json:"due_day"`
	EndDate      *time.Time `json:"end_date"`
	DiscountMode *string    `json:"discount_mode" gorm:"column:discount_mode"` // 'CLT' | 'custom' | NULL
	HourlyRate   *float64   `json:"hourly_rate" gorm:"column:hourly_rate"`     // Para calcular extras
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
