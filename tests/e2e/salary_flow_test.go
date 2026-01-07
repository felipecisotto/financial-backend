package e2e

import (
	"context"
	"testing"
	"time"

	"financial-backend/internal/domain/entities"
	"financial-backend/internal/dto"
	"financial-backend/internal/dto/request"
	"financial-backend/internal/gateways"
	incomeRepo "financial-backend/internal/repositories/income"
	salaryRepo "financial-backend/internal/repositories/salary_entry"
	incomeUC "financial-backend/internal/usecases/income"
	salaryUC "financial-backend/internal/usecases/salary"
	"financial-backend/pkg/salary"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SalaryFlowTestSuite struct {
	suite.Suite
	db              *gorm.DB
	incomeUseCase   incomeUC.UseCase
	salaryUseCase   salaryUC.UseCase
	incomeGateway   gateways.IncomeGateway
	salaryGateway   gateways.SalaryEntryGateway
	calculator      salary.Calculator
	ctx             context.Context
}

func (s *SalaryFlowTestSuite) SetupSuite() {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(s.T(), err)
	s.db = db

	// Auto migrate
	err = db.AutoMigrate(&entities.Income{}, &entities.SalaryEntry{})
	assert.NoError(s.T(), err)

	// Setup repositories
	incomeRepository := incomeRepo.NewRepository(db)
	salaryRepository := salaryRepo.NewRepository(db)

	// Setup gateways
	s.incomeGateway = gateways.NewIncomeGateway(incomeRepository)
	s.salaryGateway = gateways.NewSalaryEntryGateway(salaryRepository)

	// Setup calculator
	s.calculator = salary.NewCalculator()

	// Setup use cases
	s.incomeUseCase = incomeUC.NewUseCase(s.incomeGateway)
	s.salaryUseCase = salaryUC.NewSalaryUseCase(s.incomeGateway, s.salaryGateway, s.calculator)

	s.ctx = context.Background()
}

func (s *SalaryFlowTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (s *SalaryFlowTestSuite) SetupTest() {
	// Clean tables before each test
	s.db.Exec("DELETE FROM salary_entries")
	s.db.Exec("DELETE FROM incomes")
}

// TestSalaryCompleteFlow tests the complete salary tracking flow
func (s *SalaryFlowTestSuite) TestSalaryCompleteFlow() {
	// 1. Create Income
	now := time.Now()
	createIncomeReq := &dto.CreateIncomeRequest{
		Description: "Software Engineer Salary",
		Amount:      10000.0,
		Type:        "fixed",
		StartDate:   now,
		DueDay:      5,
	}

	incomeResp, err := s.incomeUseCase.Create(s.ctx, createIncomeReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), incomeResp)
	assert.Equal(s.T(), "SOFTWARE ENGINEER SALARY", incomeResp.Description) // Description is uppercased in model
	assert.Equal(s.T(), 10000.0, incomeResp.Amount)

	incomeID := incomeResp.ID

	// 2. Activate salary tracking (CLT)
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "CLT",
		HourlyRate:   50.0,
	}

	err = s.salaryUseCase.ActivateTracking(s.ctx, incomeID, activateReq)
	assert.NoError(s.T(), err)

	// 3. Verify INSS and IRRF entries created
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()

	summary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), summary)
	assert.Equal(s.T(), 10000.0, summary.BaseAmount)
	assert.True(s.T(), summary.TotalDiscounts > 0, "Should have CLT discounts")

	// Verify INSS and IRRF entries exist
	hasINSS := false
	hasIRRF := false
	for _, entry := range summary.Entries {
		if entry.Category == "INSS" {
			hasINSS = true
			assert.Equal(s.T(), "discount", entry.EntryType)
		}
		if entry.Category == "IRRF" {
			hasIRRF = true
			assert.Equal(s.T(), "discount", entry.EntryType)
		}
	}
	assert.True(s.T(), hasINSS, "Should have INSS entry")
	assert.True(s.T(), hasIRRF, "Should have IRRF entry")

	// 5. Add extra hours
	rate := 75.0
	addExtraHoursReq := request.AddExtraHoursRequest{
		Month: currentMonth,
		Year:  currentYear,
		Hours: 20.0,
		Rate:  &rate, // 50% extra
	}

	extraHoursSummary, err := s.salaryUseCase.AddExtraHours(s.ctx, incomeID, addExtraHoursReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), extraHoursSummary)

	// Extra hours amount should be hours * rate = 20 * 75 = 1500
	expectedExtraHours := 1500.0
	hasExtraHours := false
	for _, entry := range extraHoursSummary.Entries {
		if entry.Category == "Extra Hours" {
			hasExtraHours = true
			assert.Equal(s.T(), expectedExtraHours, entry.Amount)
			assert.Equal(s.T(), "addition", entry.EntryType)
		}
	}
	assert.True(s.T(), hasExtraHours, "Should have Extra Hours entry")

	// 6. Verify CLT entries recalculated with new gross (10000 + 1500 = 11500)
	// INSS and IRRF should be recalculated based on 11500
	expectedGrossAdjusted := 11500.0
	assert.Equal(s.T(), expectedGrossAdjusted, extraHoursSummary.BaseAmount+extraHoursSummary.TotalAdditions)

	// 7. Add on-call
	addOnCallReq := request.AddOnCallRequest{
		Month: currentMonth,
		Year:  currentYear,
		Hours: 24.0,
	}

	onCallSummary, err := s.salaryUseCase.AddOnCall(s.ctx, incomeID, addOnCallReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), onCallSummary)

	// On-call rate should be hourly_rate / 3 = 50 / 3 = 16.67
	// Amount should be hours * rate = 24 * 16.67 = 400.08 (approximately)
	expectedOnCallRate := 50.0 / 3.0
	expectedOnCallAmount := 24.0 * expectedOnCallRate

	hasOnCall := false
	for _, entry := range onCallSummary.Entries {
		if entry.Category == "On-Call" {
			hasOnCall = true
			assert.InDelta(s.T(), expectedOnCallAmount, entry.Amount, 0.01)
			assert.Equal(s.T(), "addition", entry.EntryType)
		}
	}
	assert.True(s.T(), hasOnCall, "Should have On-Call entry")

	// 8. Verify summary updated with all entries
	finalSummary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), finalSummary)

	// Should have: base gross + extra hours + on-call
	totalAdditions := expectedExtraHours + expectedOnCallAmount
	assert.InDelta(s.T(), totalAdditions, finalSummary.TotalAdditions, 0.01)

	// Net amount should be: gross + additions - discounts
	expectedNetAmount := 10000.0 + totalAdditions - finalSummary.TotalDiscounts
	assert.InDelta(s.T(), expectedNetAmount, finalSummary.NetAmount, 0.01)

	// 9. Get history for the year
	history, err := s.salaryUseCase.GetHistory(s.ctx, incomeID, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), history)
	assert.Len(s.T(), history.MonthlySummaries, 12, "Should return 12 months")

	// Verify current month has the correct net amount
	currentMonthHistory := history.MonthlySummaries[currentMonth-1]
	assert.Equal(s.T(), currentMonth, currentMonthHistory.Month)
	assert.InDelta(s.T(), expectedNetAmount, currentMonthHistory.NetAmount, 0.01)

	// 10. Remove an entry (let's remove extra hours)
	var extraHoursEntryID string
	for _, entry := range finalSummary.Entries {
		if entry.Category == "Extra Hours" {
			extraHoursEntryID = entry.ID
			break
		}
	}
	assert.NotEmpty(s.T(), extraHoursEntryID, "Extra hours entry should exist")

	err = s.salaryUseCase.RemoveEntry(s.ctx, extraHoursEntryID)
	assert.NoError(s.T(), err)

	// 11. Verify summary updated after removal
	afterRemovalSummary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), afterRemovalSummary)

	// Should only have on-call addition now
	assert.InDelta(s.T(), expectedOnCallAmount, afterRemovalSummary.TotalAdditions, 0.01)

	// Extra hours should not be in entries
	hasExtraHoursAfterRemoval := false
	for _, entry := range afterRemovalSummary.Entries {
		if entry.Category == "Extra Hours" {
			hasExtraHoursAfterRemoval = true
		}
	}
	assert.False(s.T(), hasExtraHoursAfterRemoval, "Extra hours entry should be removed")

	// Net amount should be recalculated
	expectedNetAfterRemoval := 10000.0 + expectedOnCallAmount - afterRemovalSummary.TotalDiscounts
	assert.InDelta(s.T(), expectedNetAfterRemoval, afterRemovalSummary.NetAmount, 0.01)
}

// TestSalaryCustomMode tests salary tracking with custom discount mode
func (s *SalaryFlowTestSuite) TestSalaryCustomMode() {
	// 1. Create Income
	now := time.Now()
	endDate := now.AddDate(1, 0, 0) // 1 year from now
	createIncomeReq := &dto.CreateIncomeRequest{
		Description: "Freelance Income",
		Amount:      8000.0,
		Type:        "variable",
		StartDate:   now,
		EndDate:     &endDate, // Variable income requires end date
		DueDay:      15,
	}

	incomeResp, err := s.incomeUseCase.Create(s.ctx, createIncomeReq)
	assert.NoError(s.T(), err)
	incomeID := incomeResp.ID

	// 2. Activate salary tracking with custom mode
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "custom",
		HourlyRate:   100.0,
	}

	err = s.salaryUseCase.ActivateTracking(s.ctx, incomeID, activateReq)
	assert.NoError(s.T(), err)

	// 3. Verify no CLT entries created
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()

	summary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), summary)
	assert.Equal(s.T(), 8000.0, summary.BaseAmount)
	assert.Equal(s.T(), 0.0, summary.TotalDiscounts, "Custom mode should not have automatic discounts")
	assert.Len(s.T(), summary.Entries, 0, "Should have no entries initially")

	// 4. Add custom discount entry
	desc := "Monthly health insurance"
	addEntryReq := request.CreateSalaryEntryRequest{
		Month:           currentMonth,
		Year:            currentYear,
		EntryType:       "discount",
		Category:        "Health Insurance",
		CalculationType: "fixed",
		BaseValue:       500.0,
		Multiplier:      1.0,
		Description:     &desc,
	}

	_, err = s.salaryUseCase.AddEntry(s.ctx, incomeID, addEntryReq)
	assert.NoError(s.T(), err)

	// 5. Verify custom entry added
	updatedSummary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 500.0, updatedSummary.TotalDiscounts)

	// Net amount should be gross - custom discount
	expectedNet := 8000.0 - 500.0
	assert.Equal(s.T(), expectedNet, updatedSummary.NetAmount)
}

// TestSalaryValidations tests various validation scenarios
func (s *SalaryFlowTestSuite) TestSalaryValidations() {
	// Create income for testing
	now := time.Now()
	createIncomeReq := &dto.CreateIncomeRequest{
		Description: "Test Income",
		Amount:      5000.0,
		Type:        "fixed",
		StartDate:   now,
		DueDay:      1,
	}

	incomeResp, err := s.incomeUseCase.Create(s.ctx, createIncomeReq)
	assert.NoError(s.T(), err)
	incomeID := incomeResp.ID

	// Activate tracking
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "CLT",
		HourlyRate:   40.0,
	}
	err = s.salaryUseCase.ActivateTracking(s.ctx, incomeID, activateReq)
	assert.NoError(s.T(), err)

	// Test removing non-existent entry
	err = s.salaryUseCase.RemoveEntry(s.ctx, "non-existent-id")
	assert.Error(s.T(), err, "Removing non-existent entry should return error")

	// Test getting summary successfully
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	summary, err := s.salaryUseCase.GetSummary(s.ctx, incomeID, currentMonth, currentYear)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), summary)
	assert.Equal(s.T(), 5000.0, summary.BaseAmount)
}

func TestSalaryFlowTestSuite(t *testing.T) {
	suite.Run(t, new(SalaryFlowTestSuite))
}

// Helper functions
func strPtr(s string) *string {
	return &s
}
