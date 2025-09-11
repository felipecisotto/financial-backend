package budget

import (
	"context"
	"fmt"
	"testing"
	"time"

	"financial-backend/internal/dto"
	"financial-backend/internal/domain/models"
	"financial-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BudgetUseCaseTestSuite struct {
	suite.Suite
	mockGateway *mocks.MockBudgetGateway
	useCase     UseCase
}

func (s *BudgetUseCaseTestSuite) SetupTest() {
	s.mockGateway = &mocks.MockBudgetGateway{}
	s.useCase = NewUseCase(s.mockGateway)
}

func (s *BudgetUseCaseTestSuite) TearDownTest() {
	s.mockGateway.AssertExpectations(s.T())
}

func (s *BudgetUseCaseTestSuite) TestCreate_Success() {
	input := dto.CreateBudgetRequest{
		Description: "Monthly Budget",
		Amount:      2000.0,
	}

	ctx := context.Background()

	// Mock the gateway call to succeed
	s.mockGateway.On("Create", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Create(ctx, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "MONTHLY BUDGET", result.Description)
	assert.Equal(s.T(), input.Amount, result.Amount)
	assert.NotEmpty(s.T(), result.ID)
	assert.Equal(s.T(), "active", result.Status)
}

func (s *BudgetUseCaseTestSuite) TestCreate_NegativeAmount() {
	input := dto.CreateBudgetRequest{
		Description: "Invalid Budget",
		Amount:      -1000.0,
	}

	ctx := context.Background()

	// Mock the gateway call to succeed (negative amounts are allowed to pass through)
	s.mockGateway.On("Create", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Create(ctx, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "INVALID BUDGET", result.Description)
	assert.Equal(s.T(), -1000.0, result.Amount)
	assert.NotEmpty(s.T(), result.ID)
	assert.Equal(s.T(), "active", result.Status)
}

func (s *BudgetUseCaseTestSuite) TestCreate_GatewayError() {
	input := dto.CreateBudgetRequest{
		Description: "Test Budget",
		Amount:      1000.0,
	}

	ctx := context.Background()

	// Mock the gateway call to fail
	s.mockGateway.On("Create", ctx, mock.Anything).Return(fmt.Errorf("gateway error"))

	_, err := s.useCase.Create(ctx, input)

	assert.Error(s.T(), err)
}

func (s *BudgetUseCaseTestSuite) TestUpdate_Success() {
	budgetID := "test-budget-id"
	endDate := time.Now().AddDate(0, 1, 0)
	updateReq := &dto.UpdateBudgetRequest{
		EndDate: endDate,
	}

	ctx := context.Background()

	// Create mock budget using the models package
	mockBudget := models.NewBudget(
		budgetID,
		1000.0,
		"Original Budget",
		nil,
	)

	expectedResponse := dto.BudgetResponse{
		ID:          budgetID,
		Description: "ORIGINAL BUDGET",
		Amount:      1000.0,
		EndDate:     &endDate,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockGateway.On("Get", ctx, budgetID).Return(mockBudget, nil)
	s.mockGateway.On("Update", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Update(ctx, budgetID, updateReq)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, result.ID)
	assert.Equal(s.T(), expectedResponse.Description, result.Description)
}

func (s *BudgetUseCaseTestSuite) TestUpdate_NotFound() {
	budgetID := "non-existent-id"
	endDate := time.Now().AddDate(0, 1, 0)
	updateReq := &dto.UpdateBudgetRequest{
		EndDate: endDate,
	}

	ctx := context.Background()

	s.mockGateway.On("Get", ctx, budgetID).Return(nil, fmt.Errorf("budget not found"))

	_, err := s.useCase.Update(ctx, budgetID, updateReq)

	assert.Error(s.T(), err)
}

func (s *BudgetUseCaseTestSuite) TestDelete_Success() {
	budgetID := "test-budget-id"
	ctx := context.Background()

	s.mockGateway.On("Delete", ctx, budgetID).Return(nil)

	err := s.useCase.Delete(ctx, budgetID)

	assert.NoError(s.T(), err)
}

func (s *BudgetUseCaseTestSuite) TestDelete_GatewayError() {
	budgetID := "test-budget-id"
	ctx := context.Background()

	s.mockGateway.On("Delete", ctx, budgetID).Return(fmt.Errorf("delete error"))

	err := s.useCase.Delete(ctx, budgetID)

	assert.Error(s.T(), err)
}

func (s *BudgetUseCaseTestSuite) TestGet_Success() {
	budgetID := "test-budget-id"
	ctx := context.Background()

	// Create mock budget using the models package
	mockBudget := models.NewBudget(
		budgetID,
		1000.0,
		"Test Budget",
		nil,
	)

	s.mockGateway.On("Get", ctx, budgetID).Return(mockBudget, nil)

	result, err := s.useCase.Get(ctx, budgetID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), budgetID, result.ID)
	assert.Equal(s.T(), "TEST BUDGET", result.Description)
	assert.Equal(s.T(), 1000.0, result.Amount)
}

func (s *BudgetUseCaseTestSuite) TestGet_NotFound() {
	budgetID := "non-existent-id"
	ctx := context.Background()

	s.mockGateway.On("Get", ctx, budgetID).Return(nil, fmt.Errorf("budget not found"))

	_, err := s.useCase.Get(ctx, budgetID)

	assert.Error(s.T(), err)
}

func (s *BudgetUseCaseTestSuite) TestList_Success() {
	ctx := context.Background()
	params := dto.BudgetListParams{
		Description: "Monthly",
		Status:      "ACTIVE",
		Page:        1,
		Limit:       10,
	}

	// Create mock budgets
	mockBudgets := []models.Budget{}
	budget1 := models.NewBudget("budget-1", 1000.0, "Budget 1", nil)
	budget2 := models.NewBudget("budget-2", 1500.0, "Budget 2", nil)
	mockBudgets = append(mockBudgets, budget1, budget2)

	s.mockGateway.On("List", ctx, "ACTIVE", "Monthly", mock.MatchedBy(func(page models.PageRequest) bool {
		return page.Page == 1 && page.Limit == 10
	})).Return(mockBudgets, int64(2), nil)

	result, err := s.useCase.List(ctx, params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int64(1), result.Page)
	assert.Equal(s.T(), int64(10), result.Limit)
	assert.Equal(s.T(), int64(1), result.TotalPages) // 2 items / 10 per page = 1 page
	assert.Len(s.T(), result.Results, 2)
}

func (s *BudgetUseCaseTestSuite) TestList_GatewayError() {
	ctx := context.Background()
	params := dto.BudgetListParams{
		Page:  1,
		Limit: 10,
	}

	s.mockGateway.On("List", ctx, "", "", mock.Anything).Return(nil, int64(0), fmt.Errorf("list error"))

	_, err := s.useCase.List(ctx, params)

	assert.Error(s.T(), err)
}

func TestBudgetUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BudgetUseCaseTestSuite))
}