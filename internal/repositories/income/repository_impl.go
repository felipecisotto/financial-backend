package income

import (
	"context"
	"fmt"
	"time"

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
	return r.db.WithContext(ctx).Delete(&entities.Income{}, id).Error
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Income, error) {
	var income entities.Income
	if err := r.db.WithContext(ctx).First(&income, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar receita: %v", err)
	}
	return &income, nil
}

func (r *repository) List(ctx context.Context) ([]*entities.Income, error) {
	var incomes []*entities.Income
	if err := r.db.WithContext(ctx).Find(&incomes).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar receitas: %v", err)
	}
	return incomes, nil
}

func (r *repository) ListByType(ctx context.Context, incomeType entities.IncomeType) ([]*entities.Income, error) {
	var incomes []*entities.Income
	if err := r.db.WithContext(ctx).
		Where("type = ?", incomeType).
		Find(&incomes).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar receitas: %v", err)
	}
	return incomes, nil
}

func (r *repository) ListByMonth(ctx context.Context, month time.Month, year int) ([]*entities.Income, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var incomes []*entities.Income
	if err := r.db.WithContext(ctx).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Find(&incomes).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar receitas: %v", err)
	}

	return incomes, nil
}
