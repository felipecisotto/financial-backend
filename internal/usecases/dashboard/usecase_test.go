package dashboard

import (
	"context"
	"fmt"
	"testing"

	"financial-backend/internal/dto/response"
	"financial-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DashboardUseCaseTestSuite struct {
	suite.Suite
	mockExpenseGateway        *mocks.MockExpenseGateway
	mockIncomeGateway         *mocks.MockIncomeGateway
	mockBudgetMovementGateway *mocks.MockBudgetMovementGateway
	useCase                   UseCase
}

func (s *DashboardUseCaseTestSuite) SetupTest() {
	s.mockExpenseGateway = &mocks.MockExpenseGateway{}
	s.mockIncomeGateway = &mocks.MockIncomeGateway{}
	s.mockBudgetMovementGateway = &mocks.MockBudgetMovementGateway{}
	s.useCase = NewDashBoardUseCase(
		s.mockExpenseGateway,
		s.mockIncomeGateway,
		s.mockBudgetMovementGateway,
	)
}

func (s *DashboardUseCaseTestSuite) TearDownTest() {
	s.mockExpenseGateway.AssertExpectations(s.T())
	s.mockIncomeGateway.AssertExpectations(s.T())
	s.mockBudgetMovementGateway.AssertExpectations(s.T())
}

func (s *DashboardUseCaseTestSuite) TestGetSummary_Success() {
	ctx := context.Background()
	month := 1
	year := 2024

	expectedIncome := 5000.0
	expectedExpense := 2000.0
	expectedRemaining := expectedIncome - expectedExpense

	expectedSummary := response.SummaryView{
		TotalIncome:    expectedIncome,
		TotalExpense:   expectedExpense,
		TotalRemaining: expectedRemaining,
	}

	s.mockIncomeGateway.On("SummaryByMonth", ctx, month, year).Return(expectedIncome, nil)
	s.mockExpenseGateway.On("SummaryByMonth", ctx, month, year).Return(expectedExpense, nil)

	result, err := s.useCase.GetSummary(ctx, month, year)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedSummary.TotalIncome, result.TotalIncome)
	assert.Equal(s.T(), expectedSummary.TotalExpense, result.TotalExpense)
	assert.Equal(s.T(), expectedSummary.TotalRemaining, result.TotalRemaining)
}

func (s *DashboardUseCaseTestSuite) TestGetSummary_IncomeGatewayError() {
	ctx := context.Background()
	month := 1
	year := 2024

	s.mockIncomeGateway.On("SummaryByMonth", ctx, month, year).Return(0.0, fmt.Errorf("income gateway error"))
	s.mockExpenseGateway.On("SummaryByMonth", ctx, month, year).Return(2000.0, nil)

	result, err := s.useCase.GetSummary(ctx, month, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "income gateway error")
	assert.Equal(s.T(), response.SummaryView{}, result)
}

func (s *DashboardUseCaseTestSuite) TestGetSummary_ExpenseGatewayError() {
	ctx := context.Background()
	month := 1
	year := 2024

	s.mockIncomeGateway.On("SummaryByMonth", ctx, month, year).Return(5000.0, nil)
	s.mockExpenseGateway.On("SummaryByMonth", ctx, month, year).Return(0.0, fmt.Errorf("expense gateway error"))

	result, err := s.useCase.GetSummary(ctx, month, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "expense gateway error")
	assert.Equal(s.T(), response.SummaryView{}, result)
}

func (s *DashboardUseCaseTestSuite) TestGetSummary_BothGatewayErrors() {
	ctx := context.Background()
	month := 1
	year := 2024

	s.mockIncomeGateway.On("SummaryByMonth", ctx, month, year).Return(0.0, fmt.Errorf("income gateway error"))
	s.mockExpenseGateway.On("SummaryByMonth", ctx, month, year).Return(0.0, fmt.Errorf("expense gateway error"))

	result, err := s.useCase.GetSummary(ctx, month, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "income gateway error")
	assert.Equal(s.T(), response.SummaryView{}, result)
}

func (s *DashboardUseCaseTestSuite) TestGetSummary_NegativeRemaining() {
	ctx := context.Background()
	month := 1
	year := 2024

	expectedIncome := 1000.0
	expectedExpense := 2000.0
	expectedRemaining := expectedIncome - expectedExpense

	expectedSummary := response.SummaryView{
		TotalIncome:    expectedIncome,
		TotalExpense:   expectedExpense,
		TotalRemaining: expectedRemaining,
	}

	s.mockIncomeGateway.On("SummaryByMonth", ctx, month, year).Return(expectedIncome, nil)
	s.mockExpenseGateway.On("SummaryByMonth", ctx, month, year).Return(expectedExpense, nil)

	result, err := s.useCase.GetSummary(ctx, month, year)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedSummary.TotalIncome, result.TotalIncome)
	assert.Equal(s.T(), expectedSummary.TotalExpense, result.TotalExpense)
	assert.Equal(s.T(), expectedSummary.TotalRemaining, result.TotalRemaining)
	assert.True(s.T(), result.TotalRemaining < 0, "TotalRemaining should be negative")
}

func (s *DashboardUseCaseTestSuite) TestSummaryBudgetUsageByMonthYear_Success() {
	ctx := context.Background()
	month := 1
	year := 2024

	expectedUtilizations := []response.SummaryBudgetUtilization{
		{
			Description: "Food Budget",
			Amount:      1000.0,
			Usage:       750.0,
		},
		{
			Description: "Transport Budget",
			Amount:      500.0,
			Usage:       300.0,
		},
	}

	s.mockBudgetMovementGateway.On("SummaryBudgetUsageByMonthYear", ctx, month, year).Return(expectedUtilizations, nil)

	result, err := s.useCase.SummaryBudgetUsageByMonthYear(ctx, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedUtilizations[0].Description, result[0].Description)
	assert.Equal(s.T(), expectedUtilizations[0].Amount, result[0].Amount)
	assert.Equal(s.T(), expectedUtilizations[0].Usage, result[0].Usage)
	assert.Equal(s.T(), expectedUtilizations[1].Description, result[1].Description)
	assert.Equal(s.T(), expectedUtilizations[1].Amount, result[1].Amount)
	assert.Equal(s.T(), expectedUtilizations[1].Usage, result[1].Usage)
}

func (s *DashboardUseCaseTestSuite) TestSummaryBudgetUsageByMonthYear_EmptyResult() {
	ctx := context.Background()
	month := 1
	year := 2024

	expectedUtilizations := []response.SummaryBudgetUtilization{}

	s.mockBudgetMovementGateway.On("SummaryBudgetUsageByMonthYear", ctx, month, year).Return(expectedUtilizations, nil)

	result, err := s.useCase.SummaryBudgetUsageByMonthYear(ctx, month, year)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func (s *DashboardUseCaseTestSuite) TestSummaryBudgetUsageByMonthYear_GatewayError() {
	ctx := context.Background()
	month := 1
	year := 2024

	s.mockBudgetMovementGateway.On("SummaryBudgetUsageByMonthYear", ctx, month, year).Return(nil, fmt.Errorf("budget movement gateway error"))

	result, err := s.useCase.SummaryBudgetUsageByMonthYear(ctx, month, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "budget movement gateway error")
	assert.Nil(s.T(), result)
}

func (s *DashboardUseCaseTestSuite) TestGetMonthlyEvolution_Success() {
	ctx := context.Background()
	year := 2024

	incomes := map[string]float64{
		"1": 5000.0, "2": 5200.0, "3": 5100.0, "4": 5300.0,
		"5": 5400.0, "6": 5500.0, "7": 5600.0, "8": 5700.0,
		"9": 5800.0, "10": 5900.0, "11": 6000.0, "12": 6100.0,
	}

	expenses := map[string]float64{
		"1": 3000.0, "2": 3100.0, "3": 2900.0, "4": 3200.0,
		"5": 3300.0, "6": 3400.0, "7": 3500.0, "8": 3600.0,
		"9": 3700.0, "10": 3800.0, "11": 3900.0, "12": 4000.0,
	}

	s.mockIncomeGateway.On("MonthlyEvolution", ctx, year).Return(incomes, nil)
	s.mockExpenseGateway.On("MonthlyEvolution", ctx, year).Return(expenses, nil)

	result, err := s.useCase.GetMonthlyEvolution(ctx, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 12)

	// Verify first month
	assert.Equal(s.T(), 1, result[0].Month)
	assert.Equal(s.T(), 5000.0, result[0].Income)
	assert.Equal(s.T(), 3000.0, result[0].Expense)

	// Verify last month
	assert.Equal(s.T(), 12, result[11].Month)
	assert.Equal(s.T(), 6100.0, result[11].Income)
	assert.Equal(s.T(), 4000.0, result[11].Expense)

	// Verify middle month
	assert.Equal(s.T(), 6, result[5].Month)
	assert.Equal(s.T(), 5500.0, result[5].Income)
	assert.Equal(s.T(), 3400.0, result[5].Expense)
}

func (s *DashboardUseCaseTestSuite) TestGetMonthlyEvolution_WithZeroValues() {
	ctx := context.Background()
	year := 2024

	incomes := map[string]float64{
		"1": 0.0, "2": 0.0, "3": 0.0, "4": 0.0,
		"5": 0.0, "6": 0.0, "7": 0.0, "8": 0.0,
		"9": 0.0, "10": 0.0, "11": 0.0, "12": 0.0,
	}

	expenses := map[string]float64{
		"1": 0.0, "2": 0.0, "3": 0.0, "4": 0.0,
		"5": 0.0, "6": 0.0, "7": 0.0, "8": 0.0,
		"9": 0.0, "10": 0.0, "11": 0.0, "12": 0.0,
	}

	s.mockIncomeGateway.On("MonthlyEvolution", ctx, year).Return(incomes, nil)
	s.mockExpenseGateway.On("MonthlyEvolution", ctx, year).Return(expenses, nil)

	result, err := s.useCase.GetMonthlyEvolution(ctx, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 12)

	for i := 0; i < 12; i++ {
		assert.Equal(s.T(), i+1, result[i].Month)
		assert.Equal(s.T(), 0.0, result[i].Income)
		assert.Equal(s.T(), 0.0, result[i].Expense)
	}
}

func (s *DashboardUseCaseTestSuite) TestGetMonthlyEvolution_IncomeGatewayError() {
	ctx := context.Background()
	year := 2024

	expenses := map[string]float64{
		"1": 3000.0, "2": 3100.0, "3": 2900.0, "4": 3200.0,
		"5": 3300.0, "6": 3400.0, "7": 3500.0, "8": 3600.0,
		"9": 3700.0, "10": 3800.0, "11": 3900.0, "12": 4000.0,
	}

	s.mockIncomeGateway.On("MonthlyEvolution", ctx, year).Return(map[string]float64(nil), fmt.Errorf("income evolution error"))
	s.mockExpenseGateway.On("MonthlyEvolution", ctx, year).Return(expenses, nil)

	result, err := s.useCase.GetMonthlyEvolution(ctx, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "income evolution error")
	assert.Nil(s.T(), result)
}

func (s *DashboardUseCaseTestSuite) TestGetMonthlyEvolution_ExpenseGatewayError() {
	ctx := context.Background()
	year := 2024

	incomes := map[string]float64{
		"1": 5000.0, "2": 5200.0, "3": 5100.0, "4": 5300.0,
		"5": 5400.0, "6": 5500.0, "7": 5600.0, "8": 5700.0,
		"9": 5800.0, "10": 5900.0, "11": 6000.0, "12": 6100.0,
	}

	s.mockIncomeGateway.On("MonthlyEvolution", ctx, year).Return(incomes, nil)
	s.mockExpenseGateway.On("MonthlyEvolution", ctx, year).Return(map[string]float64(nil), fmt.Errorf("expense evolution error"))

	result, err := s.useCase.GetMonthlyEvolution(ctx, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "expense evolution error")
	assert.Nil(s.T(), result)
}

func (s *DashboardUseCaseTestSuite) TestGetMonthlyEvolution_BothGatewaysError() {
	ctx := context.Background()
	year := 2024

	s.mockIncomeGateway.On("MonthlyEvolution", ctx, year).Return(map[string]float64(nil), fmt.Errorf("income evolution error"))
	s.mockExpenseGateway.On("MonthlyEvolution", ctx, year).Return(map[string]float64(nil), fmt.Errorf("expense evolution error"))

	result, err := s.useCase.GetMonthlyEvolution(ctx, year)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "income evolution error")
	assert.Nil(s.T(), result)
}

func TestDashboardUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardUseCaseTestSuite))
}
