package expense

import (
	"context"

	"financial-backend/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, expense *entities.Expense) error
	Update(ctx context.Context, expense *entities.Expense) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*entities.Expense, error)
	List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string) ([]*entities.Expense, int64, error)
}
