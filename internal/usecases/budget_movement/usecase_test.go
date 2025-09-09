package budgetmovement

import (
	"context"
	"fmt"
	"testing"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"financial-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BudgetMovementUseCaseTestSuite struct {
	suite.Suite
	mockBudgetMovementGateway *mocks.MockBudgetMovementGateway
	mockBudgetGateway         *mocks.MockBudgetGateway
	mockExpenseGateway        *mocks.MockExpenseGateway
	useCase                   UseCase
}

func (s *BudgetMovementUseCaseTestSuite) SetupTest() {
	s.mockBudgetMovementGateway = &mocks.MockBudgetMovementGateway{}
	s.mockBudgetGateway = &mocks.MockBudgetGateway{}
	s.mockExpenseGateway = &mocks.MockExpenseGateway{}
	s.useCase = NewBudgetMovementUseCase(
		s.mockBudgetMovementGateway,
		s.mockBudgetGateway,
		s.mockExpenseGateway,
	)
}

func (s *BudgetMovementUseCaseTestSuite) TearDownTest() {
	s.mockBudgetMovementGateway.AssertExpectations(s.T())
	s.mockBudgetGateway.AssertExpectations(s.T())
	s.mockExpenseGateway.AssertExpectations(s.T())
}

func (s *BudgetMovementUseCaseTestSuite) TestCreate_Success() {
	request := dtos.BudgetMovementRequest{
		BudgetId: "budget-id",
		Origin:   "MANUAL",
		Month:    1,
		Year:     2024,
		Type:     "DEBIT",
		Amount:   100,
	}

	ctx := context.Background()

	expectedResponse := dtos.BudgetMovementResponse{
		ID:     "movement-id",
		Origin: "MANUAL",
		Month:  1,
		Year:   2024,
		Type:   "DEBIT",
		Amount: 100,
		Budget: dtos.BudgetResponse{
			ID:          "budget-id",
			Description: "Test Budget",
		},
		CreatedAt: time.Now(),
	}

	// Mock the gateway call to succeed
	s.mockBudgetMovementGateway.On("Create", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Create(ctx, request)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), request.Origin, result.Origin)
	assert.Equal(s.T(), request.Month, result.Month)
	assert.Equal(s.T(), request.Year, result.Year)
	assert.Equal(s.T(), request.Type, result.Type)
	assert.Equal(s.T(), request.Amount, result.Amount)
	assert.NotEmpty(s.T(), result.ID)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreate_GatewayError() {
	request := dtos.BudgetMovementRequest{
		BudgetId: "budget-id",
		Origin:   "MANUAL",
		Month:    1,
		Year:     2024,
		Type:     "DEBIT",
		Amount:   100,
	}

	ctx := context.Background()

	// Mock the gateway call to fail
	s.mockBudgetMovementGateway.On("Create", ctx, mock.Anything).Return(fmt.Errorf("gateway error"))

	result, err := s.useCase.Create(ctx, request)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), dtos.BudgetMovementResponse{}, result)
}

func (s *BudgetMovementUseCaseTestSuite) TestFind_Success() {
	params := dtos.BudgetMovementParams{
		BudgetId:     "budget-id",
		MovementType: "DEBIT",
		Origin:       "EXPENSE",
		Month:        1,
		Year:         2024,
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	ctx := context.Background()

	// Create mock budget movements
	mockMovements := []models.BudgetMovement{}

	expectedPage := models.Page[dtos.BudgetMovementResponse]{
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Results: []dtos.BudgetMovementResponse{
			{
				ID:     "movement-1",
				Origin: "EXPENSE",
				Month:  1,
				Year:   2024,
				Type:   "DEBIT",
				Amount: 150,
			},
		},
	}

	s.mockBudgetMovementGateway.On("List", ctx, "budget-id", "DEBIT", "EXPENSE", 1, 2024, mock.MatchedBy(func(page models.PageRequest) bool {
		return page.Page == 1 && page.Limit == 10
	})).Return(mockMovements, int64(1), nil)

	result, err := s.useCase.Find(ctx, params)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), result.Page)
	assert.Equal(s.T(), int64(10), result.Limit)
}

func (s *BudgetMovementUseCaseTestSuite) TestFind_GatewayError() {
	params := dtos.BudgetMovementParams{
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	ctx := context.Background()

	s.mockBudgetMovementGateway.On("List", ctx, "", "", "", 0, 0, mock.Anything).Return(nil, int64(0), fmt.Errorf("find error"))

	result, err := s.useCase.Find(ctx, params)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), models.Page[dtos.BudgetMovementResponse]{}, result)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreateExpenseMovement_Success() {
	// Create a mock expense
	mockExpense, _ := models.NewExpense(
		"expense-id",
		"Test Expense",
		100.0,
		"VARIABLE",
		nil,
		nil,
		"PIX",
		&[]int{1}[0],
		15,
		time.Now(),
		nil,
		nil,
	)

	ctx := context.Background()

	s.mockBudgetMovementGateway.On("Create", ctx, mock.Anything).Return(nil)

	err := s.useCase.CreateExpenseMovement(ctx, mockExpense)

	assert.NoError(s.T(), err)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreateExpenseMovement_GatewayError() {
	// Create a mock expense
	mockExpense, _ := models.NewExpense(
		"expense-id",
		"Test Expense",
		100.0,
		"VARIABLE",
		nil,
		nil,
		"PIX",
		&[]int{1}[0],
		15,
		time.Now(),
		nil,
		nil,
	)

	ctx := context.Background()

	s.mockBudgetMovementGateway.On("Create", ctx, mock.Anything).Return(fmt.Errorf("create movement error"))

	err := s.useCase.CreateExpenseMovement(ctx, mockExpense)

	assert.Error(s.T(), err)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreateRecurrencyMovements_Success() {
	ctx := context.Background()

	// Mock getting expenses without movement
	mockExpenses := []models.Expense{}
	s.mockExpenseGateway.On("GetExpensesWithoutMovementInMonth", ctx).Return(mockExpenses, nil)

	err := s.useCase.CreateRecurrencyMovements(ctx)

	assert.NoError(s.T(), err)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreateRecurrencyMovements_WithExpenses() {
	ctx := context.Background()

	// Create mock expenses
	expense1, _ := models.NewExpense(
		"expense-1",
		"Recurring Expense 1",
		100.0,
		"VARIABLE",
		nil,
		nil,
		"PIX",
		&[]int{1}[0],
		15,
		time.Now(),
		nil,
		nil,
	)

	mockExpenses := []models.Expense{expense1}

	s.mockExpenseGateway.On("GetExpensesWithoutMovementInMonth", ctx).Return(mockExpenses, nil)
	s.mockBudgetMovementGateway.On("CreateAll", ctx, mock.Anything).Return(nil)

	err := s.useCase.CreateRecurrencyMovements(ctx)

	assert.NoError(s.T(), err)
}

func (s *BudgetMovementUseCaseTestSuite) TestCreateRecurrencyMovements_GatewayError() {
	ctx := context.Background()

	s.mockExpenseGateway.On("GetExpensesWithoutMovementInMonth", ctx).Return(nil, fmt.Errorf("get expenses error"))

	err := s.useCase.CreateRecurrencyMovements(ctx)

	assert.Error(s.T(), err)
}

func TestBudgetMovementUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BudgetMovementUseCaseTestSuite))
}