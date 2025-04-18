package budgetmovement

import (
	"context"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, budgetMovement entities.BudgetMovement) error
	List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) ([]entities.BudgetMovement, int64, error)
	GetById(ctx context.Context, id string) (*entities.BudgetMovement, error)
}
