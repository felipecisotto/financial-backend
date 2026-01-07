package models

import (
	"errors"
	"fmt"
	"time"
)

// Entry type constants
type EntryType string

const (
	EntryTypeDiscount EntryType = "discount"
	EntryTypeAddition EntryType = "addition"
)

// Calculation type constants
type CalculationType string

const (
	CalculationTypeFixed      CalculationType = "fixed"
	CalculationTypePercentage CalculationType = "percentage"
	CalculationTypeComputed   CalculationType = "computed"
)

// Category constants
const (
	// Discount categories
	CategoryINSS             = "INSS"
	CategoryIRRF             = "IRRF"
	CategoryHealthInsurance  = "health_insurance"
	CategoryUnionFee         = "union_fee"
	CategoryMealVoucher      = "meal_voucher"
	CategoryTransportVoucher = "transport_voucher"
	CategoryLoan             = "loan"

	// Addition categories
	CategoryExtraHours = "extra_hours"
	CategoryOnCall     = "on_call"
	CategoryBonus      = "bonus"
	CategoryCommission = "commission"
	CategoryOther      = "other"
)

// SalaryEntry represents the domain model for salary entries
type SalaryEntry interface {
	ID() string
	IncomeID() string
	Month() int
	Year() int
	EntryType() string
	Category() string
	CalculationType() string
	BaseValue() float64
	Multiplier() float64
	Amount() float64
	Description() *string
	ReferenceDate() *time.Time
	CreatedAt() time.Time
	UpdatedAt() time.Time
	IsDiscount() bool
	IsAddition() bool
}

type salaryEntry struct {
	id              string
	incomeID        string
	month           int
	year            int
	entryType       string
	category        string
	calculationType string
	baseValue       float64
	multiplier      float64
	amount          float64
	description     *string
	referenceDate   *time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

// NewSalaryEntry creates a new salary entry with validations
func NewSalaryEntry(
	id string,
	incomeID string,
	month int,
	year int,
	entryType string,
	category string,
	calculationType string,
	baseValue float64,
	multiplier float64,
	amount float64,
	description *string,
	referenceDate *time.Time,
) (SalaryEntry, error) {
	now := time.Now()

	// Validate month
	if month < 1 || month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}

	// Validate year
	if year < 2020 {
		return nil, errors.New("year must be 2020 or later")
	}

	// Validate entry type
	if entryType != string(EntryTypeDiscount) && entryType != string(EntryTypeAddition) {
		return nil, fmt.Errorf("entry_type must be '%s' or '%s'", EntryTypeDiscount, EntryTypeAddition)
	}

	// Validate calculation type
	if calculationType != string(CalculationTypeFixed) &&
		calculationType != string(CalculationTypePercentage) &&
		calculationType != string(CalculationTypeComputed) {
		return nil, fmt.Errorf("calculation_type must be '%s', '%s' or '%s'",
			CalculationTypeFixed, CalculationTypePercentage, CalculationTypeComputed)
	}

	// Validate category is not empty
	if category == "" {
		return nil, errors.New("category is required")
	}

	return &salaryEntry{
		id:              id,
		incomeID:        incomeID,
		month:           month,
		year:            year,
		entryType:       entryType,
		category:        category,
		calculationType: calculationType,
		baseValue:       baseValue,
		multiplier:      multiplier,
		amount:          amount,
		description:     description,
		referenceDate:   referenceDate,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

// Getters
func (s *salaryEntry) ID() string                  { return s.id }
func (s *salaryEntry) IncomeID() string            { return s.incomeID }
func (s *salaryEntry) Month() int                  { return s.month }
func (s *salaryEntry) Year() int                   { return s.year }
func (s *salaryEntry) EntryType() string           { return s.entryType }
func (s *salaryEntry) Category() string            { return s.category }
func (s *salaryEntry) CalculationType() string     { return s.calculationType }
func (s *salaryEntry) BaseValue() float64          { return s.baseValue }
func (s *salaryEntry) Multiplier() float64         { return s.multiplier }
func (s *salaryEntry) Amount() float64             { return s.amount }
func (s *salaryEntry) Description() *string        { return s.description }
func (s *salaryEntry) ReferenceDate() *time.Time   { return s.referenceDate }
func (s *salaryEntry) CreatedAt() time.Time        { return s.createdAt }
func (s *salaryEntry) UpdatedAt() time.Time        { return s.updatedAt }

// Helper methods
func (s *salaryEntry) IsDiscount() bool {
	return s.entryType == string(EntryTypeDiscount)
}

func (s *salaryEntry) IsAddition() bool {
	return s.entryType == string(EntryTypeAddition)
}
