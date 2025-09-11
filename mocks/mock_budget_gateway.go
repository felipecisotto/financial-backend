package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"financial-backend/internal/domain/models"
)

// MockBudgetGateway is a mock for gateways.BudgetGateway
type MockBudgetGateway struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockBudgetGateway) Create(ctx context.Context, budget models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockBudgetGateway) Update(ctx context.Context, budget models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockBudgetGateway) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockBudgetGateway) Get(ctx context.Context, id string) (models.Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.Budget), args.Error(1)
}

// List mocks the List method
func (m *MockBudgetGateway) List(ctx context.Context, status, description string, page models.PageRequest) ([]models.Budget, int64, error) {
	args := m.Called(ctx, status, description, page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Budget), args.Get(1).(int64), args.Error(2)
}

// GetBudgetsWithoutMovement mocks the GetBudgetsWithoutMovement method
func (m *MockBudgetGateway) GetBudgetsWithoutMovement(ctx context.Context) ([]models.Budget, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Budget), args.Error(1)
}
