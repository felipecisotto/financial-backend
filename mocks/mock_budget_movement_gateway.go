package mocks

import (
	"context"

	"financial-backend/internal/models"
	"financial-backend/internal/views"
	"github.com/stretchr/testify/mock"
)

// MockBudgetMovementGateway is a mock for gateways.BudgetMovementGateway
type MockBudgetMovementGateway struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockBudgetMovementGateway) Create(ctx context.Context, budgetMovement models.BudgetMovement) error {
	args := m.Called(ctx, budgetMovement)
	return args.Error(0)
}

// CreateAll mocks the CreateAll method
func (m *MockBudgetMovementGateway) CreateAll(ctx context.Context, movements []models.BudgetMovement) error {
	args := m.Called(ctx, movements)
	return args.Error(0)
}

// List mocks the List method
func (m *MockBudgetMovementGateway) List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) ([]models.BudgetMovement, int64, error) {
	args := m.Called(ctx, budgetId, movementType, origin, month, year, page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.BudgetMovement), args.Get(1).(int64), args.Error(2)
}

// GetByID mocks the GetByID method
func (m *MockBudgetMovementGateway) GetByID(ctx context.Context, id string) (models.BudgetMovement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.BudgetMovement), args.Error(1)
}

// SummaryBudgetUsageByMonthYear mocks the SummaryBudgetUsageByMonthYear method
func (m *MockBudgetMovementGateway) SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) ([]views.SummaryBudgetUtilization, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]views.SummaryBudgetUtilization), args.Error(1)
}