package entities

import (
	"time"
)

// Budget representa a tabela de orçamentos
type Budget struct {
	ID          string     `gorm:"primaryKey"`
	Description string     `gorm:"not null"`
	Amount      float64    `gorm:"not null"`
	EndDate     *time.Time `gorm:"null"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	Expenses    []Expense  `gorm:"foreignKey:BudgetID"`
}
