package entities

import (
	"time"
)

// IncomeType define os tipos de receita suportados
type IncomeType string

const (
	IncomeTypeFixed    IncomeType = "fixed"
	IncomeTypeVariable IncomeType = "variable"
)

// Income representa a tabela de receitas
type Income struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	Type        IncomeType `json:"type"`
	StartDate   time.Time
	DueDay      *int       `json:"due_day,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
