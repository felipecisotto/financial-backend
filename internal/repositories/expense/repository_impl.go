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

func (r *repository) List(ctx context.Context) ([]*entities.Expense, error) {
	var expenses []*entities.Expense
	if err := r.db.WithContext(ctx).Find(&expenses).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar despesas: %v", err)
	}
	return expenses, nil
}
