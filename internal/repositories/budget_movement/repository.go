package budgetmovement

import (
	"context"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/views"
)

type Repository interface {
	CreateAll(ctx context.Context, budgetMovements []entities.BudgetMovement) error
	Create(ctx context.Context, budgetMovement entities.BudgetMovement) error
	List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) ([]entities.BudgetMovement, int64, error)
	GetById(ctx context.Context, id string) (*entities.BudgetMovement, error)
	SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []views.SummaryBudgetUtilization, err error)
}
