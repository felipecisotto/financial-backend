package mocks

import (
	"context"

	"financial-backend/internal/domain/models"
	"financial-backend/internal/repositories/salary_entry"
	"github.com/stretchr/testify/mock"
)

// MockSalaryEntryGateway is a mock for gateways.SalaryEntryGateway
type MockSalaryEntryGateway struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockSalaryEntryGateway) Create(ctx context.Context, entry models.SalaryEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

// CreateBatch mocks the CreateBatch method
func (m *MockSalaryEntryGateway) CreateBatch(ctx context.Context, entries []models.SalaryEntry) error {
	args := m.Called(ctx, entries)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockSalaryEntryGateway) Update(ctx context.Context, entry models.SalaryEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockSalaryEntryGateway) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockSalaryEntryGateway) FindByID(ctx context.Context, id string) (models.SalaryEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.SalaryEntry), args.Error(1)
}

// FindByIncomeAndPeriod mocks the FindByIncomeAndPeriod method
func (m *MockSalaryEntryGateway) FindByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) ([]models.SalaryEntry, error) {
	args := m.Called(ctx, incomeID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SalaryEntry), args.Error(1)
}

// FindByIncomeAndCategory mocks the FindByIncomeAndCategory method
func (m *MockSalaryEntryGateway) FindByIncomeAndCategory(ctx context.Context, incomeID string, month, year int, category string) (models.SalaryEntry, error) {
	args := m.Called(ctx, incomeID, month, year, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(models.SalaryEntry), args.Error(1)
}

// DeleteByIncomeAndPeriod mocks the DeleteByIncomeAndPeriod method
func (m *MockSalaryEntryGateway) DeleteByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) error {
	args := m.Called(ctx, incomeID, month, year)
	return args.Error(0)
}

// GetMonthlySummary mocks the GetMonthlySummary method
func (m *MockSalaryEntryGateway) GetMonthlySummary(ctx context.Context, incomeID string, month, year int) (*salary_entry.MonthlySummary, error) {
	args := m.Called(ctx, incomeID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*salary_entry.MonthlySummary), args.Error(1)
}
