package income

import (
	. "context"
	. "financial-backend/internal/entities"
)

// Repository defines the interface for income repository operations
type Repository interface {
	// Create creates a new income record
	Create(ctx Context, income *Income) error

	// Update updates an existing income record
	Update(ctx Context, income *Income) error

	// Delete removes an income record by ID
	Delete(ctx Context, id string) error

	// Get retrieves an income record by ID
	Get(ctx Context, id string) (*Income, error)

	// List retrieves all income records
	List(ctx Context, incomeType, description string, limit, offset int) ([]*Income, int64, error)

	SummaryByMonth(ctx Context, month, year int) (amount float64, err error)
}
