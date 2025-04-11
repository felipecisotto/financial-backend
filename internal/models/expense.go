package models

import (
	"time"
)

// ExpenseType define os tipos de despesa suportados
type ExpenseType string

const (
	ExpenseTypeRecurring  ExpenseType = "recurring"
	ExpenseTypeCreditCard ExpenseType = "credit_card"
)

type Expense interface {
	Id() string
	Description() string
	Amount() float64
	Type() ExpenseType
	StartDate() time.Time
	DueDate() time.Time
}

// Expense representa o modelo de domínio de despesa com suas regras de negócio
type expense struct {
	id          string
	description string
	amount      float64
	expenseType ExpenseType
	startDate   time.Time
	dueDate 	time.Time
}

func NewExpense(id, description string, amount float64, expenseType string, startDate, dueDate time.Time) Expense {
	return &expense{
		id:          id,
		description: description,
		amount:      amount,
		expenseType: ExpenseType(expenseType),
		startDate: startDate,
		dueDate: dueDate,
	}
}

func (e *expense) Id() string {
	return e.id
}

func (e *expense) Description() string {
	return e.description
}

func (e *expense) Amount() float64 {
	return e.amount
}

func (e *expense) Type() ExpenseType {
	return e.expenseType
}

func (e *expense) DueDate() time.Time {
	return e.dueDate
}

func (e *expense) StartDate() time.Time {
	return e.startDate
}