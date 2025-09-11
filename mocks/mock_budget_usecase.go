package mocks

import (
	"context"

	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockBudgetUseCase is a mock for budget.UseCase
type MockBudgetUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockBudgetUseCase) Create(ctx context.Context, req request.CreateBudgetRequest) (response.BudgetResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(response.BudgetResponse), args.Error(1)
}

// Update mocks the Update method
func (m *MockBudgetUseCase) Update(ctx context.Context, id string, req *request.UpdateBudgetRequest) (response.BudgetResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(response.BudgetResponse), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockBudgetUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockBudgetUseCase) Get(ctx context.Context, id string) (response.BudgetResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(response.BudgetResponse), args.Error(1)
}

// List mocks the List method
func (m *MockBudgetUseCase) List(ctx context.Context, params request.BudgetListParams) (*models.Page[response.BudgetResponse], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page[response.BudgetResponse]), args.Error(1)
}
