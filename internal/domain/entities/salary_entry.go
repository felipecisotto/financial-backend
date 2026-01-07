package entities

import "time"

// SalaryEntry representa a tabela de lançamentos de salário
type SalaryEntry struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	IncomeID        string     `json:"income_id" gorm:"column:income_id;not null;index"`
	Month           int        `json:"month" gorm:"not null"`
	Year            int        `json:"year" gorm:"not null"`
	EntryType       string     `json:"entry_type" gorm:"not null"`
	Category        string     `json:"category" gorm:"not null"`
	CalculationType string     `json:"calculation_type" gorm:"not null"`
	BaseValue       float64    `json:"base_value" gorm:"not null;default:0"`
	Multiplier      float64    `json:"multiplier" gorm:"not null;default:1.0"`
	Amount          float64    `json:"amount" gorm:"not null"`
	Description     *string    `json:"description,omitempty"`
	ReferenceDate   *time.Time `json:"reference_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"not null"`
}

func (SalaryEntry) TableName() string {
	return "salary_entries"
}
