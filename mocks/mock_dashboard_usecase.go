package mocks

import (
	"context"

	"financial-backend/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockDashboardUseCase is a mock for dashboard.UseCase
type MockDashboardUseCase struct {
	mock.Mock
}

// GetSummary mocks the GetSummary method
func (m *MockDashboardUseCase) GetSummary(ctx context.Context, month, year int) (response.SummaryView, error) {
	args := m.Called(ctx, month, year)
	return args.Get(0).(response.SummaryView), args.Error(1)
}

// SummaryBudgetUsageByMonthYear mocks the SummaryBudgetUsageByMonthYear method
func (m *MockDashboardUseCase) SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) ([]response.SummaryBudgetUtilization, error) {
	args := m.Called(ctx, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]response.SummaryBudgetUtilization), args.Error(1)
}
