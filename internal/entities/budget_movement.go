package entities

import (
	"time"
)

type BudgetMovement struct {
	ID        string `gorm:"primaryKey"`
	BudgetId  string `gorm:"not null"`
	Budget    Budget `gorm:"foreignKey:BudgetId"`
	Origin    string
	Month     int
	Year      int
	Type      string
	Amount    int
	CreatedAt time.Time

	// field for read
	OriginDescription *string `gorm:"->;-:migration"`
}
