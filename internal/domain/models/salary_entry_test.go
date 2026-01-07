package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SalaryEntryTestSuite struct {
	suite.Suite
}

func TestSalaryEntryTestSuite(t *testing.T) {
	suite.Run(t, new(SalaryEntryTestSuite))
}

// TestNewSalaryEntry_ValidEntryCreation tests creating a valid salary entry
func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidEntryCreation() {
	description := "Test description"
	referenceDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		1000.0,
		0.11,
		110.0,
		&description,
		&referenceDate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), "test-id", entry.ID())
	assert.Equal(s.T(), "income-id", entry.IncomeID())
	assert.Equal(s.T(), 6, entry.Month())
	assert.Equal(s.T(), 2024, entry.Year())
	assert.Equal(s.T(), string(EntryTypeDiscount), entry.EntryType())
	assert.Equal(s.T(), CategoryINSS, entry.Category())
	assert.Equal(s.T(), string(CalculationTypePercentage), entry.CalculationType())
	assert.Equal(s.T(), 1000.0, entry.BaseValue())
	assert.Equal(s.T(), 0.11, entry.Multiplier())
	assert.Equal(s.T(), 110.0, entry.Amount())
	assert.Equal(s.T(), &description, entry.Description())
	assert.Equal(s.T(), &referenceDate, entry.ReferenceDate())
	assert.NotZero(s.T(), entry.CreatedAt())
	assert.NotZero(s.T(), entry.UpdatedAt())
}

// TestNewSalaryEntry_ValidEntryWithoutOptionalFields tests creating entry without optional fields
func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidEntryWithoutOptionalFields() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		12,
		2025,
		string(EntryTypeAddition),
		CategoryBonus,
		string(CalculationTypeFixed),
		500.0,
		1.0,
		500.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Nil(s.T(), entry.Description())
	assert.Nil(s.T(), entry.ReferenceDate())
}

// Month validation tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_MonthTooLow() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		0,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "month must be between 1 and 12", err.Error())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_MonthTooHigh() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		13,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "month must be between 1 and 12", err.Error())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_MonthMinimumValid() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		1,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 1, entry.Month())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_MonthMaximumValid() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		12,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 12, entry.Month())
}

// Year validation tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_YearTooLow() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2019,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "year must be 2020 or later", err.Error())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_YearMinimumValid() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2020,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 2020, entry.Year())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_YearFutureValid() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2030,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 2030, entry.Year())
}

// Entry type validation tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_InvalidEntryType() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		"invalid_type",
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "entry_type must be")
	assert.Contains(s.T(), err.Error(), string(EntryTypeDiscount))
	assert.Contains(s.T(), err.Error(), string(EntryTypeAddition))
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidDiscountType() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), string(EntryTypeDiscount), entry.EntryType())
	assert.True(s.T(), entry.IsDiscount())
	assert.False(s.T(), entry.IsAddition())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidAdditionType() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeAddition),
		CategoryBonus,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), string(EntryTypeAddition), entry.EntryType())
	assert.True(s.T(), entry.IsAddition())
	assert.False(s.T(), entry.IsDiscount())
}

// Calculation type validation tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_InvalidCalculationType() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		"invalid_calculation",
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "calculation_type must be")
	assert.Contains(s.T(), err.Error(), string(CalculationTypeFixed))
	assert.Contains(s.T(), err.Error(), string(CalculationTypePercentage))
	assert.Contains(s.T(), err.Error(), string(CalculationTypeComputed))
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidCalculationTypeFixed() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryHealthInsurance,
		string(CalculationTypeFixed),
		200.0,
		1.0,
		200.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), string(CalculationTypeFixed), entry.CalculationType())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidCalculationTypePercentage() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		1000.0,
		0.11,
		110.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), string(CalculationTypePercentage), entry.CalculationType())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidCalculationTypeComputed() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryIRRF,
		string(CalculationTypeComputed),
		5000.0,
		0.275,
		1375.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), string(CalculationTypeComputed), entry.CalculationType())
}

// Category validation tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_EmptyCategory() {
	_, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		"",
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "category is required", err.Error())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidDiscountCategories() {
	discountCategories := []string{
		CategoryINSS,
		CategoryIRRF,
		CategoryHealthInsurance,
		CategoryUnionFee,
		CategoryMealVoucher,
		CategoryTransportVoucher,
		CategoryLoan,
	}

	for _, category := range discountCategories {
		entry, err := NewSalaryEntry(
			"test-id",
			"income-id",
			6,
			2024,
			string(EntryTypeDiscount),
			category,
			string(CalculationTypeFixed),
			100.0,
			1.0,
			100.0,
			nil,
			nil,
		)

		assert.NoError(s.T(), err, "Category %s should be valid", category)
		assert.NotNil(s.T(), entry)
		assert.Equal(s.T(), category, entry.Category())
	}
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ValidAdditionCategories() {
	additionCategories := []string{
		CategoryExtraHours,
		CategoryOnCall,
		CategoryBonus,
		CategoryCommission,
		CategoryOther,
	}

	for _, category := range additionCategories {
		entry, err := NewSalaryEntry(
			"test-id",
			"income-id",
			6,
			2024,
			string(EntryTypeAddition),
			category,
			string(CalculationTypeFixed),
			100.0,
			1.0,
			100.0,
			nil,
			nil,
		)

		assert.NoError(s.T(), err, "Category %s should be valid", category)
		assert.NotNil(s.T(), entry)
		assert.Equal(s.T(), category, entry.Category())
	}
}

// Base value and multiplier tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ZeroBaseValue() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		0.0,
		1.0,
		0.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 0.0, entry.BaseValue())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_NegativeBaseValue() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		-100.0,
		1.0,
		-100.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), -100.0, entry.BaseValue())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ZeroMultiplier() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		1000.0,
		0.0,
		0.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 0.0, entry.Multiplier())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_PercentageMultiplier() {
	entry, err := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		5000.0,
		0.14,
		700.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), 0.14, entry.Multiplier())
	assert.Equal(s.T(), 700.0, entry.Amount())
}

// Helper methods tests

func (s *SalaryEntryTestSuite) TestSalaryEntry_IsDiscount() {
	discountEntry, _ := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypeFixed),
		100.0,
		1.0,
		100.0,
		nil,
		nil,
	)

	assert.True(s.T(), discountEntry.IsDiscount())
	assert.False(s.T(), discountEntry.IsAddition())
}

func (s *SalaryEntryTestSuite) TestSalaryEntry_IsAddition() {
	additionEntry, _ := NewSalaryEntry(
		"test-id",
		"income-id",
		6,
		2024,
		string(EntryTypeAddition),
		CategoryBonus,
		string(CalculationTypeFixed),
		500.0,
		1.0,
		500.0,
		nil,
		nil,
	)

	assert.True(s.T(), additionEntry.IsAddition())
	assert.False(s.T(), additionEntry.IsDiscount())
}

// Getter methods tests

func (s *SalaryEntryTestSuite) TestSalaryEntry_AllGetters() {
	description := "Monthly INSS discount"
	referenceDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	entry, err := NewSalaryEntry(
		"entry-123",
		"income-456",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		5000.0,
		0.14,
		700.0,
		&description,
		&referenceDate,
	)

	assert.NoError(s.T(), err)

	// Test all getters
	assert.Equal(s.T(), "entry-123", entry.ID())
	assert.Equal(s.T(), "income-456", entry.IncomeID())
	assert.Equal(s.T(), 6, entry.Month())
	assert.Equal(s.T(), 2024, entry.Year())
	assert.Equal(s.T(), string(EntryTypeDiscount), entry.EntryType())
	assert.Equal(s.T(), CategoryINSS, entry.Category())
	assert.Equal(s.T(), string(CalculationTypePercentage), entry.CalculationType())
	assert.Equal(s.T(), 5000.0, entry.BaseValue())
	assert.Equal(s.T(), 0.14, entry.Multiplier())
	assert.Equal(s.T(), 700.0, entry.Amount())
	assert.NotNil(s.T(), entry.Description())
	assert.Equal(s.T(), description, *entry.Description())
	assert.NotNil(s.T(), entry.ReferenceDate())
	assert.Equal(s.T(), referenceDate, *entry.ReferenceDate())
	assert.NotZero(s.T(), entry.CreatedAt())
	assert.NotZero(s.T(), entry.UpdatedAt())
	assert.True(s.T(), entry.CreatedAt().Equal(entry.UpdatedAt()))
}

// Real-world scenario tests

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_INSSDiscountScenario() {
	// INSS discount: 14% of 5000 = 700
	entry, err := NewSalaryEntry(
		"inss-entry-1",
		"salary-income",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryINSS,
		string(CalculationTypePercentage),
		5000.0,
		0.14,
		700.0,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), CategoryINSS, entry.Category())
	assert.True(s.T(), entry.IsDiscount())
	assert.Equal(s.T(), string(CalculationTypePercentage), entry.CalculationType())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_ExtraHoursScenario() {
	// Extra hours: 10 hours * 50.0 per hour = 500
	description := "10 extra hours"
	entry, err := NewSalaryEntry(
		"extra-hours-1",
		"salary-income",
		6,
		2024,
		string(EntryTypeAddition),
		CategoryExtraHours,
		string(CalculationTypeFixed),
		10.0,
		50.0,
		500.0,
		&description,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), CategoryExtraHours, entry.Category())
	assert.True(s.T(), entry.IsAddition())
	assert.Equal(s.T(), 10.0, entry.BaseValue())
	assert.Equal(s.T(), 50.0, entry.Multiplier())
	assert.Equal(s.T(), 500.0, entry.Amount())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_HealthInsuranceFixedScenario() {
	// Fixed health insurance discount
	description := "Monthly health insurance"
	entry, err := NewSalaryEntry(
		"health-insurance-1",
		"salary-income",
		6,
		2024,
		string(EntryTypeDiscount),
		CategoryHealthInsurance,
		string(CalculationTypeFixed),
		200.0,
		1.0,
		200.0,
		&description,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), CategoryHealthInsurance, entry.Category())
	assert.True(s.T(), entry.IsDiscount())
	assert.Equal(s.T(), string(CalculationTypeFixed), entry.CalculationType())
	assert.Equal(s.T(), 200.0, entry.Amount())
}

func (s *SalaryEntryTestSuite) TestNewSalaryEntry_BonusScenario() {
	// Performance bonus
	description := "Q2 Performance Bonus"
	referenceDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	entry, err := NewSalaryEntry(
		"bonus-1",
		"salary-income",
		6,
		2024,
		string(EntryTypeAddition),
		CategoryBonus,
		string(CalculationTypeFixed),
		2000.0,
		1.0,
		2000.0,
		&description,
		&referenceDate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), entry)
	assert.Equal(s.T(), CategoryBonus, entry.Category())
	assert.True(s.T(), entry.IsAddition())
	assert.Equal(s.T(), 2000.0, entry.Amount())
	assert.NotNil(s.T(), entry.ReferenceDate())
}
