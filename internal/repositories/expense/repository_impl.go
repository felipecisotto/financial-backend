package expense

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

func (r *repository) Create(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *repository) Update(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Expense{}, id).Error
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

func (r *repository) ListByMonth(ctx context.Context, month time.Month, year int) ([]*entities.Expense, error) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	var expenses []*entities.Expense
	if err := r.db.WithContext(ctx).
		Where("(type = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)) OR "+
			"(type = ? AND due_date BETWEEN ? AND ?)",
			models.ExpenseTypeRecurring, endDate, startDate,
			models.ExpenseTypeCreditCard, startDate, endDate).
		Find(&expenses).Error; err != nil {
		return nil, fmt.Errorf("erro ao listar despesas: %v", err)
	}

	return expenses, nil
}

func (r *repository) AssignToBudget(ctx context.Context, expenseID string, budgetID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Expense{}).
		Where("id = ?", expenseID).
		Update("budget_id", budgetID).Error
}

func (r *repository) RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Expense{}).
		Where("id = ? AND budget_id = ?", expenseID, budgetID).
		Update("budget_id", nil).Error
}