package mocks

import (
	"context"

	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto"
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
func (m *MockBudgetMovementUseCase) Create(ctx context.Context, request dto.BudgetMovementRequest) (dto.BudgetMovementResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(dto.BudgetMovementResponse), args.Error(1)
}

// Find mocks the Find method
func (m *MockBudgetMovementUseCase) Find(ctx context.Context, params dto.BudgetMovementParams) (models.Page[dto.BudgetMovementResponse], error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Page[dto.BudgetMovementResponse]), args.Error(1)
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
