package expense

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

type ExpenseUseCaseTestSuite struct {
	suite.Suite
	mockExpenseGateway *mocks.MockExpenseGateway
	mockBudgetGateway  *mocks.MockBudgetGateway
	mockPublisher      *mocks.MockPublisher
	useCase            UseCase
}

func (s *ExpenseUseCaseTestSuite) SetupTest() {
	s.mockExpenseGateway = &mocks.MockExpenseGateway{}
	s.mockBudgetGateway = &mocks.MockBudgetGateway{}
	s.mockPublisher = &mocks.MockPublisher{}
	
	s.useCase = NewUseCase(s.mockExpenseGateway, s.mockBudgetGateway, s.mockPublisher, 15)
}

func (s *ExpenseUseCaseTestSuite) TearDownTest() {
	s.mockExpenseGateway.AssertExpectations(s.T())
	s.mockBudgetGateway.AssertExpectations(s.T())
	s.mockPublisher.AssertExpectations(s.T())
}

func (s *ExpenseUseCaseTestSuite) TestCreate_Success() {
	input := &dtos.ExpenseDTO{
		Description:  "Test Expense",
		Amount:       100.0,
		Type:         "VARIABLE",
		Method:       "PIX",
		DueDay:       15,
		StartDate:    time.Now(),
		Installments: &[]int{1}[0],
	}

	ctx := context.Background()

	// Mock the gateway call to succeed
	s.mockExpenseGateway.On("Create", ctx, mock.Anything).Return(nil)
	
	// Mock the publisher call
	s.mockPublisher.On("Publish", mock.Anything).Return()

	result, err := s.useCase.Create(ctx, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), input.Description, result.Description)
	assert.Equal(s.T(), input.Amount, result.Amount)
	assert.Equal(s.T(), input.Type, result.Type)
	assert.NotEmpty(s.T(), result.ID)
}

func (s *ExpenseUseCaseTestSuite) TestCreate_GatewayError() {
	input := &dtos.ExpenseDTO{
		Description:  "Test Expense",
		Amount:       100.0,
		Type:         "VARIABLE",
		Method:       "PIX",
		DueDay:       15,
		StartDate:    time.Now(),
		Installments: &[]int{1}[0],
	}

	ctx := context.Background()

	// Mock the gateway call to fail
	s.mockExpenseGateway.On("Create", ctx, mock.Anything).Return(fmt.Errorf("gateway error"))

	result, err := s.useCase.Create(ctx, input)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "erro ao criar despesa")
}

func (s *ExpenseUseCaseTestSuite) TestCreate_CreditCardExpense() {
	input := &dtos.ExpenseDTO{
		Description:  "Credit Card Expense",
		Amount:       500.0,
		Type:         "VARIABLE",
		Method:       string(models.ExpenseMethodCreditCard),
		DueDay:       15,
		StartDate:    time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), // Day 10
		Installments: &[]int{3}[0],
	}

	ctx := context.Background()

	// Mock the gateway call to succeed
	s.mockExpenseGateway.On("Create", ctx, mock.MatchedBy(func(exp models.Expense) bool {
		// For credit card, start date should be adjusted to due date
		expectedStartDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Due date
		return exp.StartDate().Day() == expectedStartDate.Day() &&
			exp.EndDate() != nil // Should have end date for installments
	})).Return(nil)
	
	// Mock the publisher call
	s.mockPublisher.On("Publish", mock.Anything).Return()

	result, err := s.useCase.Create(ctx, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), input.Description, result.Description)
	assert.Equal(s.T(), input.Amount, result.Amount)
}

func (s *ExpenseUseCaseTestSuite) TestDelete_Success() {
	expenseID := "test-expense-id"
	ctx := context.Background()

	s.mockExpenseGateway.On("Delete", ctx, expenseID).Return(nil)

	err := s.useCase.Delete(ctx, expenseID)

	assert.NoError(s.T(), err)
}

func (s *ExpenseUseCaseTestSuite) TestDelete_GatewayError() {
	expenseID := "test-expense-id"
	ctx := context.Background()

	s.mockExpenseGateway.On("Delete", ctx, expenseID).Return(fmt.Errorf("gateway error"))

	err := s.useCase.Delete(ctx, expenseID)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "erro ao excluir despesa")
}

func (s *ExpenseUseCaseTestSuite) TestFindByID_Success() {
	expenseID := "test-expense-id"
	ctx := context.Background()

	// Create a mock expense using the existing models.NewExpense function
	mockExpense, _ := models.NewExpense(
		expenseID,
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

	s.mockExpenseGateway.On("Get", ctx, expenseID).Return(mockExpense, nil)

	result, err := s.useCase.FindByID(ctx, expenseID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expenseID, result.ID)
	assert.Equal(s.T(), "Test Expense", result.Description)
}

func (s *ExpenseUseCaseTestSuite) TestFindByID_NotFound() {
	expenseID := "non-existent-id"
	ctx := context.Background()

	s.mockExpenseGateway.On("Get", ctx, expenseID).Return(nil, fmt.Errorf("expense not found"))

	result, err := s.useCase.FindByID(ctx, expenseID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "erro ao buscar despesa")
}

func (s *ExpenseUseCaseTestSuite) TestList_Success() {
	ctx := context.Background()
	request := &dtos.ListExpensesRequest{
		Description: "Test",
		Type:        "VARIABLE",
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	// Create mock expenses
	mockExpenses := []models.Expense{}
	expense1, _ := models.NewExpense(
		"expense-1",
		"Expense 1",
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
	expense2, _ := models.NewExpense(
		"expense-2",
		"Expense 2",
		200.0,
		"VARIABLE",
		nil,
		nil,
		"PIX",
		&[]int{1}[0],
		20,
		time.Now(),
		nil,
		nil,
	)
	mockExpenses = append(mockExpenses, expense1, expense2)

	s.mockExpenseGateway.On("List", ctx, "Test", "VARIABLE", "", "", "", "", mock.MatchedBy(func(page models.PageRequest) bool {
		return page.Page == 1 && page.Limit == 10
	})).Return(mockExpenses, int64(2), nil)

	result, err := s.useCase.List(ctx, request)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int64(1), result.Page)
	assert.Equal(s.T(), int64(10), result.Limit)
	assert.Equal(s.T(), int64(1), result.TotalPages) // 2 items / 10 per page = 1 page
	assert.Len(s.T(), result.Results, 2)
	assert.Equal(s.T(), "expense-1", result.Results[0].ID)
	assert.Equal(s.T(), "expense-2", result.Results[1].ID)
}

func (s *ExpenseUseCaseTestSuite) TestList_GatewayError() {
	ctx := context.Background()
	request := &dtos.ListExpensesRequest{
		Description: "Test",
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	s.mockExpenseGateway.On("List", ctx, "Test", "", "", "", "", "", mock.Anything).Return(nil, int64(0), fmt.Errorf("gateway error"))

	result, err := s.useCase.List(ctx, request)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "erro ao listar despesas")
}

func TestExpenseUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseUseCaseTestSuite))
}