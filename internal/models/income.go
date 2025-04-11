package models

import (
	"time"
)

type IncomeType string

const (
	IncomeTypeFixed    IncomeType = "fixed"
	IncomeTypeVariable IncomeType = "variable"
)

type Income interface {
	ID() string
	Description() string
	Amount() float64
	Type() IncomeType
	DueDay() *int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type income struct {
	id          string
	description string
	amount      float64
	incomeType  IncomeType
	dueDay      *int
	createdAt   time.Time
	updatedAt   time.Time
}

func NewIncome(id, description string, amount float64, incomeType IncomeType, dueDay *int) Income {
	now := time.Now()
	return &income{
		id:          id,
		description: description,
		amount:      amount,
		incomeType:  incomeType,
		dueDay: dueDay,
		createdAt:   now,
		updatedAt:   now,
	}
}

func (i *income) ID() string {
	return i.id
}

func (i *income) Description() string {
	return i.description
}

func (i *income) Amount() float64 {
	return i.amount
}

func (i *income) Type() IncomeType {
	return i.incomeType
}

func (i *income) DueDay() *int {
	return i.dueDay
}

func (i *income) CreatedAt() time.Time {
	return i.createdAt
}

func (i *income) UpdatedAt() time.Time {
	return i.updatedAt
}

