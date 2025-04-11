package budget

import (
	"context"
	"fmt"
	"time"

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

func (r *repository) List(ctx context.Context, page models.PageRequest) (budgets []entities.Budget, count int64, err error) {
	if err = r.db.WithContext(ctx).Offset(page.Offset()).Limit(int(page.Limit)).Find(&budgets).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar orçamentos: %v", err)
	}
	if err := r.db.WithContext(ctx).Model(&budgets).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar orçamentos: %v", err)
	}
	return budgets, count, nil
}

func (r *repository) ListByMonth(ctx context.Context, month time.Month, year int) ([]*entities.Budget, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var budgets []*entities.Budget
	if err := r.db.WithContext(ctx).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Find(&budgets).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar orçamentos: %v", err)
	}

	return budgets, nil
}

func (r *repository) GetSummary(ctx context.Context, id string) (*entities.Budget, error) {
	var budget entities.Budget
	if err := r.db.WithContext(ctx).
		Preload("Expenses").
		First(&budget, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar orçamento: %v", err)
	}

	return &budget, nil
}

func (r *repository) CreateBudgetExpense(ctx context.Context, budgetID string, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verifica se o orçamento existe
		var budget entities.Budget
		if err := tx.First(&budget, "id = ?", budgetID).Error; err != nil {
			return err
		}

		// Atualiza o ID do orçamento na despesa
		expense.BudgetID = &budgetID

		// Cria a despesa
		if err := tx.Create(expense).Error; err != nil {
			return err
		}

		return nil
	})
}
