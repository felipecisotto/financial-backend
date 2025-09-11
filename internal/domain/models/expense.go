package models

import (
	"fmt"
	"time"
)

type ExpenseType string
type ExpenseRecurrency string
type ExpenseMethod string

const (
	ExpenseTypeRecurring ExpenseType = "recurring"
	ExpenseTypeSingle    ExpenseType = "single"
)

const (
	ExpenseRecurrencyMonthly ExpenseRecurrency = "monthly"
	ExpenseRecurrencyWeekly  ExpenseRecurrency = "weekly"
	ExpenseRecurrencyDaily   ExpenseRecurrency = "daily"
)

const (
	ExpenseMethodCreditCard ExpenseMethod = "credit_card"
	ExpenseMethodPix        ExpenseMethod = "pix"
	ExpenseMethodBankSlip   ExpenseMethod = "bank_slip"
)

type Expense interface {
	Id() string
	Description() string
	Amount() float64
	Type() ExpenseType
	Recurrency() *ExpenseRecurrency
	Method() ExpenseMethod
	Installments() *int
	DueDay() int
	BudgetId() *string
	Budget() *Budget
	StartDate() time.Time
	EndDate() *time.Time
}

// Expense representa o modelo de domínio de despesa com suas regras de negócio
type expense struct {
	id           string
	description  string
	amount       float64
	expenseType  ExpenseType
	recurrency   *ExpenseRecurrency
	method       ExpenseMethod
	installments *int
	dueDay       int
	budgetId     *string
	budget       *Budget
	startDate    time.Time
	endDate      *time.Time
}

func NewExpense(
	id, description string,
	amount float64,
	expenseType string,
	budgetId,
	recurrency *string,
	method string,
	installments *int,
	dueDay int,
	startDate time.Time,
	endDate *time.Time,
	budget *Budget,
) (Expense, error) {
	if expenseType == string(ExpenseTypeRecurring) && recurrency == nil {
		return nil, fmt.Errorf("quando o tipo de despesa é recorrente, é necessário ter preencher a recorencia")
	}

	var expenseRecurrency *ExpenseRecurrency

	if recurrency == nil {
		expenseRecurrency = nil
	} else {
		recurrencyValue := *recurrency
		expenseRecurrency = new(ExpenseRecurrency)              // Cria uma nova instância
		*expenseRecurrency = ExpenseRecurrency(recurrencyValue) // Atribui o valor
	}

	return &expense{
		id:           id,
		description:  description,
		amount:       amount,
		expenseType:  ExpenseType(expenseType),
		recurrency:   expenseRecurrency,
		method:       ExpenseMethod(method),
		installments: installments,
		dueDay:       dueDay,
		budgetId:     budgetId,
		startDate:    startDate,
		endDate:      endDate,
		budget:       budget,
	}, nil
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

func (e *expense) Recurrency() *ExpenseRecurrency {
	return e.recurrency
}

func (e *expense) Method() ExpenseMethod {
	return e.method
}

func (e *expense) Installments() *int {
	return e.installments
}

func (e *expense) StartDate() time.Time {
	return e.startDate
}

func (e *expense) Budget() *Budget {
	return e.budget
}

func (e *expense) BudgetId() *string {
	return e.budgetId
}

func (e *expense) DueDay() int {
	return e.dueDay
}

func (e *expense) EndDate() *time.Time {
	return e.endDate
}
