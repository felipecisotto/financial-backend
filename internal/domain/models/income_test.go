package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IncomeTestSuite struct {
	suite.Suite
}

func TestIncomeTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeTestSuite))
}

// TestNewIncome_ValidIncomeCreation tests creating a valid income
func (s *IncomeTestSuite) TestNewIncome_ValidIncomeCreation() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	income, err := NewIncome(
		"test-id",
		"Monthly Salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "test-id", income.ID())
	assert.Equal(s.T(), "MONTHLY SALARY", income.Description()) // Description is uppercased
	assert.Equal(s.T(), 5000.0, income.Amount())
	assert.Equal(s.T(), IncomeTypeFixed, income.Type())
	assert.Equal(s.T(), 15, income.DueDay())
	assert.Equal(s.T(), startDate, income.StartDate())
	assert.Nil(s.T(), income.EndDate())
	assert.Nil(s.T(), income.DiscountMode())
	assert.Nil(s.T(), income.HourlyRate())
	assert.False(s.T(), income.HasSalaryTracking())
	assert.NotZero(s.T(), income.CreatedAt())
	assert.NotZero(s.T(), income.UpdatedAt())
}

// TestNewIncome_DescriptionUppercase tests that description is converted to uppercase
func (s *IncomeTestSuite) TestNewIncome_DescriptionUppercase() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"monthly salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "MONTHLY SALARY", income.Description())
}

// TestNewIncome_VariableWithEndDate tests creating variable income with end date
func (s *IncomeTestSuite) TestNewIncome_VariableWithEndDate() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	income, err := NewIncome(
		"test-id",
		"Temporary Project",
		3000.0,
		IncomeTypeVariable,
		10,
		startDate,
		&endDate,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), IncomeTypeVariable, income.Type())
	assert.NotNil(s.T(), income.EndDate())
	assert.Equal(s.T(), endDate, *income.EndDate())
}

// TestNewIncome_VariableWithoutEndDate tests that variable income requires end date
func (s *IncomeTestSuite) TestNewIncome_VariableWithoutEndDate() {
	startDate := time.Now()

	_, err := NewIncome(
		"test-id",
		"Variable Income",
		3000.0,
		IncomeTypeVariable,
		10,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "receita váriavel é obrigatório data final", err.Error())
}

// TestNewIncome_FixedWithoutEndDate tests that fixed income doesn't require end date
func (s *IncomeTestSuite) TestNewIncome_FixedWithoutEndDate() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Fixed Salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Nil(s.T(), income.EndDate())
}

// DiscountMode validation tests

func (s *IncomeTestSuite) TestNewIncome_ValidDiscountModeCLT() {
	startDate := time.Now()
	discountMode := "CLT"

	income, err := NewIncome(
		"test-id",
		"CLT Salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.NotNil(s.T(), income.DiscountMode())
	assert.Equal(s.T(), "CLT", *income.DiscountMode())
	assert.True(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_ValidDiscountModeCustom() {
	startDate := time.Now()
	discountMode := "custom"

	income, err := NewIncome(
		"test-id",
		"Custom Salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.NotNil(s.T(), income.DiscountMode())
	assert.Equal(s.T(), "custom", *income.DiscountMode())
	assert.True(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_InvalidDiscountMode() {
	startDate := time.Now()
	discountMode := "PJ"

	_, err := NewIncome(
		"test-id",
		"Invalid Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		nil,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "discount_mode deve ser 'CLT' ou 'custom'", err.Error())
}

func (s *IncomeTestSuite) TestNewIncome_DiscountModeNil() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Simple Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Nil(s.T(), income.DiscountMode())
	assert.False(s.T(), income.HasSalaryTracking())
}

// HourlyRate validation tests

func (s *IncomeTestSuite) TestNewIncome_ValidHourlyRate() {
	startDate := time.Now()
	hourlyRate := 50.0

	income, err := NewIncome(
		"test-id",
		"Hourly Work",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.NotNil(s.T(), income.HourlyRate())
	assert.Equal(s.T(), 50.0, *income.HourlyRate())
}

func (s *IncomeTestSuite) TestNewIncome_HourlyRateZero() {
	startDate := time.Now()
	hourlyRate := 0.0

	_, err := NewIncome(
		"test-id",
		"Invalid Hourly",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		&hourlyRate,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "hourly_rate deve ser maior que zero", err.Error())
}

func (s *IncomeTestSuite) TestNewIncome_HourlyRateNegative() {
	startDate := time.Now()
	hourlyRate := -50.0

	_, err := NewIncome(
		"test-id",
		"Invalid Hourly",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		&hourlyRate,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "hourly_rate deve ser maior que zero", err.Error())
}

func (s *IncomeTestSuite) TestNewIncome_HourlyRateNil() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Simple Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Nil(s.T(), income.HourlyRate())
}

// HasSalaryTracking tests

func (s *IncomeTestSuite) TestNewIncome_HasSalaryTrackingWithDiscountMode() {
	startDate := time.Now()
	discountMode := "CLT"

	income, err := NewIncome(
		"test-id",
		"CLT Salary",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.True(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_HasSalaryTrackingWithoutDiscountMode() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Simple Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.False(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_HasSalaryTrackingWithHourlyRateOnly() {
	startDate := time.Now()
	hourlyRate := 50.0

	// Having only hourly rate without discount mode means no salary tracking
	income, err := NewIncome(
		"test-id",
		"Hourly Work",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.False(s.T(), income.HasSalaryTracking())
}

// Combined validation tests

func (s *IncomeTestSuite) TestNewIncome_WithAllSalaryTrackingFields() {
	startDate := time.Now()
	discountMode := "custom"
	hourlyRate := 75.5

	income, err := NewIncome(
		"test-id",
		"Full Salary Tracking",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "custom", *income.DiscountMode())
	assert.Equal(s.T(), 75.5, *income.HourlyRate())
	assert.True(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_InvalidDiscountModeWithValidHourlyRate() {
	startDate := time.Now()
	discountMode := "invalid"
	hourlyRate := 50.0

	_, err := NewIncome(
		"test-id",
		"Invalid Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		&hourlyRate,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "discount_mode deve ser 'CLT' ou 'custom'", err.Error())
}

func (s *IncomeTestSuite) TestNewIncome_ValidDiscountModeWithInvalidHourlyRate() {
	startDate := time.Now()
	discountMode := "CLT"
	hourlyRate := 0.0

	_, err := NewIncome(
		"test-id",
		"Invalid Income",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		&discountMode,
		&hourlyRate,
	)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "hourly_rate deve ser maior que zero", err.Error())
}

// Getter methods tests

func (s *IncomeTestSuite) TestNewIncome_AllGetters() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	discountMode := "CLT"
	hourlyRate := 62.5

	income, err := NewIncome(
		"income-123",
		"main salary",
		6000.0,
		IncomeTypeFixed,
		20,
		startDate,
		&endDate,
		&discountMode,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)

	// Test all getters
	assert.Equal(s.T(), "income-123", income.ID())
	assert.Equal(s.T(), "MAIN SALARY", income.Description())
	assert.Equal(s.T(), 6000.0, income.Amount())
	assert.Equal(s.T(), IncomeTypeFixed, income.Type())
	assert.Equal(s.T(), 20, income.DueDay())
	assert.Equal(s.T(), startDate, income.StartDate())
	assert.NotNil(s.T(), income.EndDate())
	assert.Equal(s.T(), endDate, *income.EndDate())
	assert.NotNil(s.T(), income.DiscountMode())
	assert.Equal(s.T(), "CLT", *income.DiscountMode())
	assert.NotNil(s.T(), income.HourlyRate())
	assert.Equal(s.T(), 62.5, *income.HourlyRate())
	assert.True(s.T(), income.HasSalaryTracking())
	assert.NotZero(s.T(), income.CreatedAt())
	assert.NotZero(s.T(), income.UpdatedAt())
	assert.True(s.T(), income.CreatedAt().Equal(income.UpdatedAt()))
}

// Real-world scenario tests

func (s *IncomeTestSuite) TestNewIncome_CLTSalaryScenario() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	discountMode := "CLT"
	hourlyRate := 50.0

	income, err := NewIncome(
		"clt-salary-1",
		"Software Engineer - CLT",
		8000.0,
		IncomeTypeFixed,
		5,
		startDate,
		nil,
		&discountMode,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "SOFTWARE ENGINEER - CLT", income.Description())
	assert.Equal(s.T(), 8000.0, income.Amount())
	assert.Equal(s.T(), IncomeTypeFixed, income.Type())
	assert.Equal(s.T(), "CLT", *income.DiscountMode())
	assert.Equal(s.T(), 50.0, *income.HourlyRate())
	assert.True(s.T(), income.HasSalaryTracking())
	assert.Nil(s.T(), income.EndDate())
}

func (s *IncomeTestSuite) TestNewIncome_FreelanceProjectScenario() {
	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC)

	income, err := NewIncome(
		"freelance-1",
		"Q3 Freelance Project",
		12000.0,
		IncomeTypeVariable,
		15,
		startDate,
		&endDate,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), IncomeTypeVariable, income.Type())
	assert.NotNil(s.T(), income.EndDate())
	assert.False(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_CustomDiscountScenario() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	discountMode := "custom"
	hourlyRate := 100.0

	income, err := NewIncome(
		"custom-1",
		"Consultant - Custom Deductions",
		15000.0,
		IncomeTypeFixed,
		10,
		startDate,
		nil,
		&discountMode,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "custom", *income.DiscountMode())
	assert.Equal(s.T(), 100.0, *income.HourlyRate())
	assert.True(s.T(), income.HasSalaryTracking())
}

func (s *IncomeTestSuite) TestNewIncome_SimpleFixedIncomeScenario() {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	income, err := NewIncome(
		"rental-1",
		"Monthly Rental Income",
		2500.0,
		IncomeTypeFixed,
		1,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "MONTHLY RENTAL INCOME", income.Description())
	assert.Equal(s.T(), IncomeTypeFixed, income.Type())
	assert.Nil(s.T(), income.DiscountMode())
	assert.Nil(s.T(), income.HourlyRate())
	assert.False(s.T(), income.HasSalaryTracking())
}

// Edge cases

func (s *IncomeTestSuite) TestNewIncome_ZeroAmount() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Zero Income",
		0.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), 0.0, income.Amount())
}

func (s *IncomeTestSuite) TestNewIncome_NegativeAmount() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Negative Income",
		-1000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), -1000.0, income.Amount())
}

func (s *IncomeTestSuite) TestNewIncome_DueDayZero() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Zero Due Day",
		5000.0,
		IncomeTypeFixed,
		0,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), 0, income.DueDay())
}

func (s *IncomeTestSuite) TestNewIncome_DueDayThirtyOne() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Last Day Income",
		5000.0,
		IncomeTypeFixed,
		31,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), 31, income.DueDay())
}

func (s *IncomeTestSuite) TestNewIncome_EmptyDescription() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), "", income.Description())
}

func (s *IncomeTestSuite) TestNewIncome_VeryLargeAmount() {
	startDate := time.Now()

	income, err := NewIncome(
		"test-id",
		"Large Income",
		999999999.99,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		nil,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), 999999999.99, income.Amount())
}

func (s *IncomeTestSuite) TestNewIncome_VeryHighHourlyRate() {
	startDate := time.Now()
	hourlyRate := 500.75

	income, err := NewIncome(
		"test-id",
		"High Rate",
		5000.0,
		IncomeTypeFixed,
		15,
		startDate,
		nil,
		nil,
		&hourlyRate,
	)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), income)
	assert.Equal(s.T(), 500.75, *income.HourlyRate())
}
