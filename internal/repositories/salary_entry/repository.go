package salary_entry

import (
	"context"
	"financial-backend/internal/domain/entities"
)

// MonthlySummary represents the aggregated summary for a month
type MonthlySummary struct {
	IncomeID       string
	Month          int
	Year           int
	BaseAmount     float64
	TotalAdditions float64
	TotalDiscounts float64
	NetAmount      float64
}

// Repository defines the interface for salary entry repository operations
type Repository interface {
	// Create creates a new salary entry record
	Create(ctx context.Context, entry *entities.SalaryEntry) error

	// CreateBatch creates multiple salary entry records in a single transaction
	CreateBatch(ctx context.Context, entries []*entities.SalaryEntry) error

	// Update updates an existing salary entry record
	Update(ctx context.Context, entry *entities.SalaryEntry) error

	// Delete removes a salary entry record by ID
	Delete(ctx context.Context, id string) error

	// FindByID retrieves a salary entry record by ID
	FindByID(ctx context.Context, id string) (*entities.SalaryEntry, error)

	// FindByIncomeAndPeriod retrieves all salary entries for a specific income and period
	// Results are ordered by: additions first, then discounts, then by created_at
	FindByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) ([]*entities.SalaryEntry, error)

	// FindByIncomeAndCategory retrieves a specific salary entry by income, period, and category
	FindByIncomeAndCategory(ctx context.Context, incomeID string, month, year int, category string) (*entities.SalaryEntry, error)

	// DeleteByIncomeAndPeriod removes all salary entries for a specific income and period
	DeleteByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) error

	// GetMonthlySummary calculates the summary for a specific income and period
	// Includes base amount (from income), total additions, total discounts, and net amount
	GetMonthlySummary(ctx context.Context, incomeID string, month, year int) (*MonthlySummary, error)
}
