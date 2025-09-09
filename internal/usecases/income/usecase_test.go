package income

import (
	"context"
	"fmt"
	"testing"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"financial-backend/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IncomeUseCaseTestSuite struct {
	suite.Suite
	mockGateway *mocks.MockIncomeGateway
	useCase     UseCase
}

func (s *IncomeUseCaseTestSuite) SetupTest() {
	s.mockGateway = &mocks.MockIncomeGateway{}
	s.useCase = NewUseCase(s.mockGateway)
}

func (s *IncomeUseCaseTestSuite) TearDownTest() {
	s.mockGateway.AssertExpectations(s.T())
}

func (s *IncomeUseCaseTestSuite) TestCreate_Success() {
	input := &dtos.CreateIncomeRequest{
		Description: "Salary",
		Amount:      5000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
	}

	ctx := context.Background()

	// Mock the gateway call to succeed
	s.mockGateway.On("Create", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Create(ctx, input)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), input.Description, result.Description)
	assert.Equal(s.T(), input.Amount, result.Amount)
	assert.Equal(s.T(), input.Type, result.Type)
	assert.NotEmpty(s.T(), result.ID)
}

func (s *IncomeUseCaseTestSuite) TestCreate_NegativeAmount() {
	input := &dtos.CreateIncomeRequest{
		Description: "Invalid Income",
		Amount:      -1000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
	}

	ctx := context.Background()

	result, err := s.useCase.Create(ctx, input)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "valor deve ser positivo")
}

func (s *IncomeUseCaseTestSuite) TestCreate_GatewayError() {
	input := &dtos.CreateIncomeRequest{
		Description: "Test Income",
		Amount:      1000.0,
		Type:        "VARIABLE",
		DueDay:      5,
		StartDate:   time.Now(),
	}

	ctx := context.Background()

	// Mock the gateway call to fail
	s.mockGateway.On("Create", ctx, mock.Anything).Return(fmt.Errorf("gateway error"))

	result, err := s.useCase.Create(ctx, input)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *IncomeUseCaseTestSuite) TestUpdate_Success() {
	incomeID := "test-income-id"
	amount := 6000.0
	updateReq := &dtos.UpdateIncomeRequest{
		Amount: &amount,
	}

	ctx := context.Background()

	// Create mock income using the models package
	mockIncome, _ := models.NewIncome(
		incomeID,
		"Original Income",
		5000.0,
		"FIXED",
		1,
		time.Now(),
		nil,
	)

	expectedResponse := &dtos.IncomeResponse{
		ID:          incomeID,
		Description: "Original Income",
		Amount:      6000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockGateway.On("Get", ctx, incomeID).Return(mockIncome, nil)
	s.mockGateway.On("Update", ctx, mock.Anything).Return(nil)

	result, err := s.useCase.Update(ctx, incomeID, updateReq)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedResponse.Amount, result.Amount)
}

func (s *IncomeUseCaseTestSuite) TestUpdate_NotFound() {
	incomeID := "non-existent-id"
	amount := 6000.0
	updateReq := &dtos.UpdateIncomeRequest{
		Amount: &amount,
	}

	ctx := context.Background()

	s.mockGateway.On("Get", ctx, incomeID).Return(nil, fmt.Errorf("income not found"))

	result, err := s.useCase.Update(ctx, incomeID, updateReq)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *IncomeUseCaseTestSuite) TestDelete_Success() {
	incomeID := "test-income-id"
	ctx := context.Background()

	s.mockGateway.On("Delete", ctx, incomeID).Return(nil)

	err := s.useCase.Delete(ctx, incomeID)

	assert.NoError(s.T(), err)
}

func (s *IncomeUseCaseTestSuite) TestDelete_GatewayError() {
	incomeID := "test-income-id"
	ctx := context.Background()

	s.mockGateway.On("Delete", ctx, incomeID).Return(fmt.Errorf("delete error"))

	err := s.useCase.Delete(ctx, incomeID)

	assert.Error(s.T(), err)
}

func (s *IncomeUseCaseTestSuite) TestGet_Success() {
	incomeID := "test-income-id"
	ctx := context.Background()

	// Create mock income using the models package
	mockIncome, _ := models.NewIncome(
		incomeID,
		"Test Income",
		1000.0,
		"FIXED",
		1,
		time.Now(),
		nil,
	)

	s.mockGateway.On("Get", ctx, incomeID).Return(mockIncome, nil)

	result, err := s.useCase.Get(ctx, incomeID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), incomeID, result.ID)
	assert.Equal(s.T(), "Test Income", result.Description)
}

func (s *IncomeUseCaseTestSuite) TestGet_NotFound() {
	incomeID := "non-existent-id"
	ctx := context.Background()

	s.mockGateway.On("Get", ctx, incomeID).Return(nil, fmt.Errorf("income not found"))

	result, err := s.useCase.Get(ctx, incomeID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *IncomeUseCaseTestSuite) TestList_Success() {
	ctx := context.Background()
	params := dtos.ListIncomeParams{
		Type:        "FIXED",
		Description: "Salary",
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	// Create mock incomes
	mockIncomes := []models.Income{}
	income1, _ := models.NewIncome("income-1", "Income 1", 1000.0, "FIXED", 1, time.Now(), nil)
	income2, _ := models.NewIncome("income-2", "Income 2", 1500.0, "FIXED", 5, time.Now(), nil)
	mockIncomes = append(mockIncomes, income1, income2)

	s.mockGateway.On("List", ctx, "FIXED", "Salary", mock.MatchedBy(func(page models.PageRequest) bool {
		return page.Page == 1 && page.Limit == 10
	})).Return(mockIncomes, int64(2), nil)

	result, err := s.useCase.List(ctx, params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int64(1), result.Page)
	assert.Equal(s.T(), int64(10), result.Limit)
	assert.Equal(s.T(), int64(1), result.TotalPages) // 2 items / 10 per page = 1 page
	assert.Len(s.T(), result.Results, 2)
}

func (s *IncomeUseCaseTestSuite) TestList_GatewayError() {
	ctx := context.Background()
	params := dtos.ListIncomeParams{
		PageRequest: dtos.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}

	s.mockGateway.On("List", ctx, "", "", mock.Anything).Return(nil, int64(0), fmt.Errorf("list error"))

	result, err := s.useCase.List(ctx, params)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func TestIncomeUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeUseCaseTestSuite))
}