package mocks

import (
	"context"

	"financial-backend/internal/dto"
	"financial-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockExpenseUseCase is a mock for expense.UseCase
type MockExpenseUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockExpenseUseCase) Create(ctx context.Context, input *dto.ExpenseDTO) (*dto.ExpenseResponse, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dto.ExpenseResponse), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockExpenseUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockExpenseUseCase) FindByID(ctx context.Context, id string) (*dto.ExpenseResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ExpenseResponse), args.Error(1)
}

// List mocks the List method
func (m *MockExpenseUseCase) List(ctx context.Context, input *dto.ListExpensesRequest) (*models.Page[*dto.ExpenseResponse], error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page[*dto.ExpenseResponse]), args.Error(1)
}