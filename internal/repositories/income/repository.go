package income

import (
	"context"

	"financial-backend/internal/entities"
)

// Repository defines the interface for income repository operations
type Repository interface {
	// Create creates a new income record
	Create(ctx context.Context, income *entities.Income) error

	// Update updates an existing income record
	Update(ctx context.Context, income *entities.Income) error

	// Delete removes an income record by ID
	Delete(ctx context.Context, id string) error

	// Get retrieves an income record by ID
	Get(ctx context.Context, id string) (*entities.Income, error)

	// List retrieves all income records
	List(ctx context.Context, incomeType, description string, limit, offset int) ([]*entities.Income, int64, error)
}
