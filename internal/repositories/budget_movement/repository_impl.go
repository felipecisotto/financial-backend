package budgetmovement

import (
	"context"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(ctx context.Context, budgetMovement entities.BudgetMovement) error {
	return r.db.WithContext(ctx).Create(budgetMovement).Error
}

// GetById implements Repository.
func (r *repository) GetById(ctx context.Context, id string) (budget *entities.BudgetMovement, err error) {
	if err := r.db.WithContext(ctx).First(&budget, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return
}

// List implements Repository.
func (r *repository) List(ctx context.Context, page models.PageRequest) (budgets []entities.BudgetMovement, count int64, err error) {
	query := r.db.WithContext(ctx)

	if err := query.Offset(page.Offset()).Limit(int(page.Limit)).Find(budgets).Error; err != nil {
		return make([]entities.BudgetMovement, 0), 0, err
	}

	if err := query.Model(&entities.BudgetMovement{}).Count(&count).Error; err != nil {
		return make([]entities.BudgetMovement, 0), 0, err
	}

	return
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
