package models

import (
	"time"
)

type BudgetStatus string

const (
	BudgetActive  BudgetStatus = "active"
	BudgetExpired BudgetStatus = "expired"
)

type Budget interface {
	ID() string
	Amount() float64
	Description() string
	Status() BudgetStatus
	EndDate() *time.Time
	CreatedAt() time.Time
	UpdatedAt() time.Time

	SetEndDate(endDate time.Time)
}

type budget struct {
	id          string
	amount      float64
	description string
	endDate     *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBudget(id string, amount float64, description string, endDate *time.Time) Budget {
	now := time.Now()
	return &budget{
		id:          id,
		amount:      amount,
		description: description,
		endDate:     endDate,
		createdAt:   now,
		updatedAt:   now,
	}
}

func (b *budget) ID() string {
	return b.id
}

func (b *budget) Amount() float64 {
	return b.amount
}

func (b *budget) Description() string {
	return b.description
}

func (b *budget) EndDate() *time.Time {
	return b.endDate
}

func (b *budget) CreatedAt() time.Time {
	return b.createdAt
}

func (b *budget) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *budget) SetEndDate(endDate time.Time) {
	b.endDate = &endDate
}

func (b *budget) Status() BudgetStatus {
	if b.endDate != nil && b.endDate.Before(time.Now()) {
		return BudgetExpired
	}
	return BudgetActive
}

type BudgetSummary interface {
	TotalBudgeted() float64
	TotalSpent() float64
	Remaining() float64
	PercentageSpent() float64
}

type budgetSummary struct {
	totalBudgeted   float64
	totalSpent      float64
	remaining       float64
	percentageSpent float64
}

func NewBudgetSummary(totalBudgeted, totalSpent float64) BudgetSummary {
	remaining := totalBudgeted - totalSpent
	percentageSpent := 0.0
	if totalBudgeted > 0 {
		percentageSpent = (totalSpent / totalBudgeted) * 100
	}

	return &budgetSummary{
		totalBudgeted:   totalBudgeted,
		totalSpent:      totalSpent,
		remaining:       remaining,
		percentageSpent: percentageSpent,
	}
}

func (s *budgetSummary) TotalBudgeted() float64 {
	return s.totalBudgeted
}

func (s *budgetSummary) TotalSpent() float64 {
	return s.totalSpent
}

func (s *budgetSummary) Remaining() float64 {
	return s.remaining
}

func (s *budgetSummary) PercentageSpent() float64 {
	return s.percentageSpent
}
