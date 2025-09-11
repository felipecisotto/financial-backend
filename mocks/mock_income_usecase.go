package mocks

import (
	"context"

	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"financial-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockIncomeUseCase is a mock for income.UseCase
type MockIncomeUseCase struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockIncomeUseCase) Create(ctx context.Context, req *request.CreateIncomeRequest) (*response.IncomeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.IncomeResponse), args.Error(1)
}

// Update mocks the Update method
func (m *MockIncomeUseCase) Update(ctx context.Context, id string, req *request.UpdateIncomeRequest) (*response.IncomeResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.IncomeResponse), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockIncomeUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockIncomeUseCase) Get(ctx context.Context, id string) (*response.IncomeResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.IncomeResponse), args.Error(1)
}

// List mocks the List method
func (m *MockIncomeUseCase) List(ctx context.Context, params request.ListIncomeParams) (*models.Page[*response.IncomeResponse], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page[*response.IncomeResponse]), args.Error(1)
}