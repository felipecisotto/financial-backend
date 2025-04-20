package budget

import (
	"context"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, budget *entities.Budget) error
	Update(ctx context.Context, budget *entities.Budget) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*entities.Budget, error)
	List(ctx context.Context, status string, description string, page models.PageRequest) ([]entities.Budget, int64, error)
	GetBudgetsWithoutMovement(ctx context.Context) ([]entities.Budget, error)
}
