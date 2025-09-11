package mocks

import (
	"context"

	"financial-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockIncomeGateway is a mock for gateways.IncomeGateway
type MockIncomeGateway struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockIncomeGateway) Create(ctx context.Context, income models.Income) error {
	args := m.Called(ctx, income)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockIncomeGateway) Update(ctx context.Context, income models.Income) error {
	args := m.Called(ctx, income)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockIncomeGateway) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Get mocks the Get method
func (m *MockIncomeGateway) Get(ctx context.Context, id string) (models.Income, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.Income), args.Error(1)
}

// List mocks the List method
func (m *MockIncomeGateway) List(ctx context.Context, incomeType, description string, page models.PageRequest) ([]models.Income, int64, error) {
	args := m.Called(ctx, incomeType, description, page)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Income), args.Get(1).(int64), args.Error(2)
}

// SummaryByMonth mocks the SummaryByMonth method
func (m *MockIncomeGateway) SummaryByMonth(ctx context.Context, month, year int) (float64, error) {
	args := m.Called(ctx, month, year)
	return args.Get(0).(float64), args.Error(1)
}
