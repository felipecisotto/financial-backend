package mocks

import (
	"context"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockBudgetMovementUseCase is a mock for budget_movement.UseCase
type MockBudgetMovementUseCase struct {
	mock.Mock
}

// FindByID mocks the FindByID method
func (m *MockBudgetMovementUseCase) FindByID(ctx context.Context, id string) {
	m.Called(ctx, id)
}

// Create mocks the Create method
func (m *MockBudgetMovementUseCase) Create(ctx context.Context, request dtos.BudgetMovementRequest) (dtos.BudgetMovementResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(dtos.BudgetMovementResponse), args.Error(1)
}

// Find mocks the Find method
func (m *MockBudgetMovementUseCase) Find(ctx context.Context, params dtos.BudgetMovementParams) (models.Page[dtos.BudgetMovementResponse], error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Page[dtos.BudgetMovementResponse]), args.Error(1)
}

// CreateExpenseMovement mocks the CreateExpenseMovement method
func (m *MockBudgetMovementUseCase) CreateExpenseMovement(ctx context.Context, expense models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

// CreateRecurrencyMovements mocks the CreateRecurrencyMovements method
func (m *MockBudgetMovementUseCase) CreateRecurrencyMovements(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}