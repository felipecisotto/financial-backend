package factories

import (
	"strings"

	"github.com/google/uuid"
)

// Factories holds all factory instances
type Factories struct {
	Expense        *ExpenseFactory
	Income         *IncomeFactory
	Budget         *BudgetFactory
	BudgetMovement *BudgetMovementFactory
}

// NewFactories creates all factories
func NewFactories() *Factories {
	return &Factories{
		Expense:        NewExpenseFactory(),
		Income:         NewIncomeFactory(),
		Budget:         NewBudgetFactory(),
		BudgetMovement: NewBudgetMovementFactory(),
	}
}

// Common helper functions used across factories
func NormalizeString(input string) string {
	return strings.TrimSpace(input)
}

func NormalizeUpperString(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

func GenerateID() string {
	return uuid.New().String()
}