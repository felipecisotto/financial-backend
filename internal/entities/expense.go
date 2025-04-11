package entities

import (
	"time"
)

// Expense representa a tabela de despesas
type Expense struct {
	ID          string  `gorm:"primaryKey"`
	Description string  `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	Type        string  `gorm:"not null"`
	BudgetID    *string `gorm:"index"`
	StartDate   time.Time
	DueDate     time.Time
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}
