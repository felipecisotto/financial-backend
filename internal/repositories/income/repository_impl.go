package income

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

func (r *repository) Create(ctx context.Context, income *entities.Income) error {
	return r.db.WithContext(ctx).Create(income).Error
}

func (r *repository) Update(ctx context.Context, income *entities.Income) error {
	return r.db.WithContext(ctx).Save(income).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Income{}).Error
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Income, error) {
	var income entities.Income
	if err := r.db.WithContext(ctx).First(&income, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar receita: %v", err)
	}
	return &income, nil
}

func (r *repository) List(ctx context.Context, incomeType, description string, limit, offset int) ([]*entities.Income, int64, error) {
	var incomes []*entities.Income
	var count int64

	query := r.db.WithContext(ctx)

	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}

	if incomeType != "" {
		query = query.Where("type = ?", incomeType)
	}

	if err := query.Offset(offset).Limit(limit).Find(&incomes).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar receitas: %v", err)
	}

	if err := query.Model(&entities.Income{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar receitas: %v", err)
	}
	return incomes, count, nil
}
