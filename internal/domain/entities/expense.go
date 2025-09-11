package entities

import (
	"time"
)

// Expense representa a tabela de despesas
type Expense struct {
	ID           string  `gorm:"primaryKey"`
	Description  string  `gorm:"not null"`
	Amount       float64 `gorm:"not null"`
	Type         string  `gorm:"not null"`
	BudgetID     *string `gorm:"index"`
	Budget       *Budget
	Recurrency   *string
	Method       string
	Installments *int
	StartDate    time.Time
	DueDay       int
	EndDate      *time.Time
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}
