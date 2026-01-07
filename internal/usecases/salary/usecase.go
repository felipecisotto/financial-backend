package salary

import (
	"context"
	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
)

// UseCase defines the interface for salary-related business operations
type UseCase interface {
	// ActivateTracking activates salary tracking for an income
	ActivateTracking(ctx context.Context, incomeID string, req request.ActivateSalaryTrackingRequest) error

	// AddEntry adds a generic salary entry
	AddEntry(ctx context.Context, incomeID string, req request.CreateSalaryEntryRequest) (*response.SalarySummaryResponse, error)

	// AddExtraHours adds extra hours to the salary
	AddExtraHours(ctx context.Context, incomeID string, req request.AddExtraHoursRequest) (*response.SalarySummaryResponse, error)

	// AddOnCall adds on-call hours to the salary
	AddOnCall(ctx context.Context, incomeID string, req request.AddOnCallRequest) (*response.SalarySummaryResponse, error)

	// RemoveEntry removes a salary entry by ID
	RemoveEntry(ctx context.Context, entryID string) error

	// GetSummary retrieves the monthly salary summary
	GetSummary(ctx context.Context, incomeID string, month, year int) (*response.SalarySummaryResponse, error)

	// GetHistory retrieves the annual salary history
	GetHistory(ctx context.Context, incomeID string, year int) (*response.SalaryHistoryResponse, error)

	// RecalculateCLT recalculates CLT entries for a specific month
	RecalculateCLT(ctx context.Context, incomeID string, month, year int) error
}
