package models

import (
	"time"
)

type MovementType string

const (
	MovementTransfer MovementType = "transfer"
	MovementIncome   MovementType = "income"
	MovementExpense  MovementType = "expense"
)

// BudgetMovementInterface defines the methods for BudgetMovement
type BudgetMovement interface {
	ID() string
	BudgetId() string
	Budget() Budget
	Origin() string
	Month() int
	Year() int
	Type() string
	Amount() int
	CreatedAt() time.Time
}

// BudgetMovement struct implements BudgetMovementInterface
type budgetMovement struct {
	id           string
	budgetId     string
	budget       Budget
	origin       string
	month        int
	year         int
	movementType string
	amount       int
	createdAt    time.Time
}

// NewBudgetMovement creates a new BudgetMovement instance
func NewBudgetMovement(
	id string,
	budgetId string,
	budget Budget,
	origin string,
	month int,
	year int,
	movementType string,
	amount int,
) BudgetMovement {
	return &budgetMovement{
		id:           id,
		budgetId:     budgetId,
		budget:       budget,
		origin:       origin,
		month:        month,
		year:         year,
		movementType: movementType,
		amount:       amount,
		createdAt:    time.Now(),
	}
}

// ID returns the ID of the BudgetMovement
func (bm *budgetMovement) ID() string {
	return bm.id
}

// BudgetId returns the Budget ID associated with the BudgetMovement
func (bm *budgetMovement) BudgetId() string {
	return bm.budgetId
}

// Budget implements BudgetMovement.
func (bm *budgetMovement) Budget() Budget {
	return bm.budget
}

// Origin implements BudgetMovement.
func (bm *budgetMovement) Origin() string {
	return bm.origin
}

// Month returns the month of the BudgetMovement
func (bm *budgetMovement) Month() int {
	return bm.month
}

// Year returns the year of the BudgetMovement
func (bm *budgetMovement) Year() int {
	return bm.year
}

// Type returns the type of the BudgetMovement
func (bm *budgetMovement) Type() string {
	return bm.movementType
}

// Amount returns the amount of the BudgetMovement
func (bm *budgetMovement) Amount() int {
	return bm.amount
}

// CreatedAt returns the creation time of the BudgetMovement
func (bm *budgetMovement) CreatedAt() time.Time {
	return bm.createdAt
}
