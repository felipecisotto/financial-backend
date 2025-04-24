package budget

import (
	"context"
	"fmt"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, budget *entities.Budget) error {
	return r.db.WithContext(ctx).Create(budget).Error
}

func (r *repository) Update(ctx context.Context, budget *entities.Budget) error {
	return r.db.WithContext(ctx).Save(budget).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Budget{}).Error
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Budget, error) {
	var budget entities.Budget
	if err := r.db.WithContext(ctx).First(&budget, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar orçamento: %v", err)
	}
	return &budget, nil
}

func (r *repository) List(ctx context.Context, status string, description string, page models.PageRequest) (budgets []entities.Budget, count int64, err error) {
	query := r.db.WithContext(ctx)

	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}

	switch status {
	case "active":
		query = query.Where("end_date is null or end_date >= CURRENT_DATE")
	case "expired":
		query = query.Where("end_date is not null and end_date < CURRENT_DATE")
	default:
	}
	if err = query.Offset(page.Offset()).Limit(int(page.Limit)).Find(&budgets).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar orçamentos: %v", err)
	}

	if err := query.Model(&entities.Budget{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar orçamentos: %v", err)
	}

	return budgets, count, nil
}

func (r *repository) GetBudgetsWithoutMovement(ctx context.Context) (reponses []entities.Budget, err error) {
	query := `select *
from budgets b
where (end_date is null or end_date >= current_date)
  and not exists(select 1
                 from budget_movements bm
                 where bm.month = cast(to_char(now() :: date, 'MM') as numeric)
                   and bm.year = cast(to_char(now() :: date, 'YYYY') as numeric)
                   and bm.type = 'start'
                   and bm.budget_id = b.id)`
	if err := r.db.WithContext(ctx).Raw(query).Find(&reponses).Error; err != nil {
		return []entities.Budget{}, err
	}
	return
}
