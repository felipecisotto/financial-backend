package mocks

import (
	"financial-backend/pkg/salary"
	"github.com/stretchr/testify/mock"
)

// MockCalculator is a mock for salary.Calculator
type MockCalculator struct {
	mock.Mock
}

// CalculateCLTDiscounts mocks the CalculateCLTDiscounts method
func (m *MockCalculator) CalculateCLTDiscounts(grossAmount float64, month, year int) ([]salary.Entry, error) {
	args := m.Called(grossAmount, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]salary.Entry), args.Error(1)
}

// CalculateExtraHours mocks the CalculateExtraHours method
func (m *MockCalculator) CalculateExtraHours(hourlyRate float64, hours float64) float64 {
	args := m.Called(hourlyRate, hours)
	return args.Get(0).(float64)
}

// CalculateOnCall mocks the CalculateOnCall method
func (m *MockCalculator) CalculateOnCall(hourlyRate float64, hours float64) float64 {
	args := m.Called(hourlyRate, hours)
	return args.Get(0).(float64)
}
