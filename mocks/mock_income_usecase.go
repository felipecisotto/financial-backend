package mocks

import (
	"context"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockIncomeUseCase is a mock for income.UseCase
type MockIncomeUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockIncomeUseCase) Create(ctx context.Context, dto *dtos.CreateIncomeRequest) (*dtos.IncomeResponse, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.IncomeResponse), args.Error(1)
}

// Update mocks the Update method
func (m *MockIncomeUseCase) Update(ctx context.Context, id string, dto *dtos.UpdateIncomeRequest) (*dtos.IncomeResponse, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.IncomeResponse), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockIncomeUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockIncomeUseCase) Get(ctx context.Context, id string) (*dtos.IncomeResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.IncomeResponse), args.Error(1)
}

// List mocks the List method
func (m *MockIncomeUseCase) List(ctx context.Context, params dtos.ListIncomeParams) (*models.Page[*dtos.IncomeResponse], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page[*dtos.IncomeResponse]), args.Error(1)
}