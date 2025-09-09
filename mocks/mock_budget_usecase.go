package mocks

import (
	"context"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockBudgetUseCase is a mock for budget.UseCase
type MockBudgetUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockBudgetUseCase) Create(ctx context.Context, dto dtos.CreateBudgetRequest) (dtos.BudgetResponse, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(dtos.BudgetResponse), args.Error(1)
}

// Update mocks the Update method
func (m *MockBudgetUseCase) Update(ctx context.Context, id string, dto *dtos.UpdateBudgetRequest) (dtos.BudgetResponse, error) {
	args := m.Called(ctx, id, dto)
	return args.Get(0).(dtos.BudgetResponse), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockBudgetUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockBudgetUseCase) Get(ctx context.Context, id string) (dtos.BudgetResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.BudgetResponse), args.Error(1)
}

// List mocks the List method
func (m *MockBudgetUseCase) List(ctx context.Context, params dtos.BudgetListParams) (*models.Page[dtos.BudgetResponse], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page[dtos.BudgetResponse]), args.Error(1)
}