package budget

import (
	"context"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, budget *entities.Budget) error
	Update(ctx context.Context, budget *entities.Budget) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*entities.Budget, error)
	List(ctx context.Context, page models.PageRequest) ([]entities.Budget, int64, error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]*entities.Budget, error)
	GetSummary(ctx context.Context, id string) (*entities.Budget, error)
}
