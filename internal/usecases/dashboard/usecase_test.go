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

func TestDashboardUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardUseCaseTestSuite))
}