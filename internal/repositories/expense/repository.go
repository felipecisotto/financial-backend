package expense

import (
	"context"
	"time"

	"financial-backend/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, expense *entities.Expense) error
	Update(ctx context.Context, expense *entities.Expense) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*entities.Expense, error)
	List(ctx context.Context) ([]*entities.Expense, error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]*entities.Expense, error)
	AssignToBudget(ctx context.Context, expenseID string, budgetID string) error
	RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error
}
