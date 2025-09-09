package mocks

import (
	"context"

	"financial-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockExpenseGateway is a mock for gateways.ExpenseGateway
type MockExpenseGateway struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockExpenseGateway) Create(ctx context.Context, expense models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockExpenseGateway) Update(ctx context.Context, expense models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockExpenseGateway) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockExpenseGateway) Get(ctx context.Context, id string) (models.Expense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.Expense), args.Error(1)
}

// List mocks the List method
func (m *MockExpenseGateway) List(ctx context.Context, description, expenseType, category, budgetId, recurrency, method string, page models.PageRequest) ([]models.Expense, int64, error) {
	args := m.Called(ctx, description, expenseType, category, budgetId, recurrency, method, page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Expense), args.Get(1).(int64), args.Error(2)
}

// GetExpensesWithoutMovementInMonth mocks the GetExpensesWithoutMovementInMonth method
func (m *MockExpenseGateway) GetExpensesWithoutMovementInMonth(ctx context.Context) ([]models.Expense, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expense), args.Error(1)
}

// SummaryByMonth mocks the SummaryByMonth method
func (m *MockExpenseGateway) SummaryByMonth(ctx context.Context, month, year int) (float64, error) {
	args := m.Called(ctx, month, year)
	return args.Get(0).(float64), args.Error(1)
}