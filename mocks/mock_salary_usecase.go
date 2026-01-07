package mocks

import (
	"context"

	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"github.com/stretchr/testify/mock"
)

// MockSalaryUseCase is a mock for salary.UseCase
type MockSalaryUseCase struct {
	mock.Mock
}

// ActivateTracking mocks the ActivateTracking method
func (m *MockSalaryUseCase) ActivateTracking(ctx context.Context, incomeID string, req request.ActivateSalaryTrackingRequest) error {
	args := m.Called(ctx, incomeID, req)
	return args.Error(0)
}

// AddEntry mocks the AddEntry method
func (m *MockSalaryUseCase) AddEntry(ctx context.Context, incomeID string, req request.CreateSalaryEntryRequest) (*response.SalarySummaryResponse, error) {
	args := m.Called(ctx, incomeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SalarySummaryResponse), args.Error(1)
}

// AddExtraHours mocks the AddExtraHours method
func (m *MockSalaryUseCase) AddExtraHours(ctx context.Context, incomeID string, req request.AddExtraHoursRequest) (*response.SalarySummaryResponse, error) {
	args := m.Called(ctx, incomeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SalarySummaryResponse), args.Error(1)
}

// AddOnCall mocks the AddOnCall method
func (m *MockSalaryUseCase) AddOnCall(ctx context.Context, incomeID string, req request.AddOnCallRequest) (*response.SalarySummaryResponse, error) {
	args := m.Called(ctx, incomeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SalarySummaryResponse), args.Error(1)
}

// RemoveEntry mocks the RemoveEntry method
func (m *MockSalaryUseCase) RemoveEntry(ctx context.Context, entryID string) error {
	args := m.Called(ctx, entryID)
	return args.Error(0)
}

// GetSummary mocks the GetSummary method
func (m *MockSalaryUseCase) GetSummary(ctx context.Context, incomeID string, month, year int) (*response.SalarySummaryResponse, error) {
	args := m.Called(ctx, incomeID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SalarySummaryResponse), args.Error(1)
}

// GetHistory mocks the GetHistory method
func (m *MockSalaryUseCase) GetHistory(ctx context.Context, incomeID string, year int) (*response.SalaryHistoryResponse, error) {
	args := m.Called(ctx, incomeID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SalaryHistoryResponse), args.Error(1)
}

// RecalculateCLT mocks the RecalculateCLT method
func (m *MockSalaryUseCase) RecalculateCLT(ctx context.Context, incomeID string, month, year int) error {
	args := m.Called(ctx, incomeID, month, year)
	return args.Error(0)
}
