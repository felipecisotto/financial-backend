package expense

import (
	"context"
	"fmt"

	"financial-backend/internal/entities"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *repository) Update(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Expense{}).Error
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Expense, error) {
	var expense entities.Expense
	if err := r.db.WithContext(ctx).First(&expense, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar despesa: %v", err)
	}
	return &expense, nil
}

func (r *repository) List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string) (expenses []*entities.Expense, count int64, err error) {
	query := r.db.WithContext(ctx)

	if description != "" {
		query = query.Where("description like ?", "%"+description+"%")
	}

	if expenseType != "" {
		query = query.Where("type = ?", expenseType)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if budgetId != "" {
		query = query.Where("budget_id = ?", budgetId)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if method != "" {
		query = query.Where("method = ?", method)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if err := query.Find(&expenses).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar despesas: %v", err)
	}

	if err := query.Model(&entities.Expense{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar Despesas: %v", err)
	}

	return expenses, count, nil
}
