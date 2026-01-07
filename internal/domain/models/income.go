package models

import (
	"errors"
	"strings"
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
	DueDay() int
	StartDate() time.Time
	EndDate() *time.Time
	DiscountMode() *string
	HourlyRate() *float64
	HasSalaryTracking() bool
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type income struct {
	id           string
	description  string
	amount       float64
	incomeType   IncomeType
	dueDay       int
	startDate    time.Time
	endDate      *time.Time
	discountMode *string
	hourlyRate   *float64
	createdAt    time.Time
	updatedAt    time.Time
}

func NewIncome(
	id, description string,
	amount float64,
	incomeType IncomeType,
	dueDay int,
	startDate time.Time,
	endDate *time.Time,
	discountMode *string,
	hourlyRate *float64,
) (Income, error) {
	now := time.Now()

	if incomeType == IncomeTypeVariable && endDate == nil {
		return nil, errors.New("receita váriavel é obrigatório data final")
	}

	// Validations for salary tracking fields
	if discountMode != nil {
		mode := *discountMode
		if mode != "CLT" && mode != "custom" {
			return nil, errors.New("discount_mode deve ser 'CLT' ou 'custom'")
		}
	}

	if hourlyRate != nil && *hourlyRate <= 0 {
		return nil, errors.New("hourly_rate deve ser maior que zero")
	}

	return &income{
		id:           id,
		description:  strings.ToUpper(description),
		amount:       amount,
		incomeType:   incomeType,
		dueDay:       dueDay,
		startDate:    startDate,
		endDate:      endDate,
		discountMode: discountMode,
		hourlyRate:   hourlyRate,
		createdAt:    now,
		updatedAt:    now,
	}, nil
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

func (i *income) DueDay() int {
	return i.dueDay
}

func (i *income) StartDate() time.Time {
	return i.startDate
}

func (i *income) EndDate() *time.Time {
	return i.endDate
}

func (i *income) CreatedAt() time.Time {
	return i.createdAt
}

func (i *income) DiscountMode() *string {
	return i.discountMode
}

func (i *income) HourlyRate() *float64 {
	return i.hourlyRate
}

func (i *income) HasSalaryTracking() bool {
	return i.discountMode != nil
}

func (i *income) UpdatedAt() time.Time {
	return i.updatedAt
}
