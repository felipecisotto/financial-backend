package expense

import (
	"context"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, expense *entities.Expense) error
	Update(ctx context.Context, expense *entities.Expense) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*entities.Expense, error)
	List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string, page models.PageRequest) ([]*entities.Expense, int64, error)
	GetExpensesWithoutMovimentInMonth(ctx context.Context) ([]*entities.Expense, error)
	SummaryByMonth(ctx context.Context, month int, year int) (float64, error)
}
