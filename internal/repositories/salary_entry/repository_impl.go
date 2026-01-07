package salary_entry

import (
	"context"
	"fmt"

	"financial-backend/internal/domain/entities"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the salary entry repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entry *entities.SalaryEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *repository) CreateBatch(ctx context.Context, entries []*entities.SalaryEntry) error {
	return r.db.WithContext(ctx).Create(&entries).Error
}

func (r *repository) Update(ctx context.Context, entry *entities.SalaryEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.SalaryEntry{}).Error
}

func (r *repository) FindByID(ctx context.Context, id string) (*entities.SalaryEntry, error) {
	var entry entities.SalaryEntry
	if err := r.db.WithContext(ctx).First(&entry, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar lançamento de salário: %v", err)
	}
	return &entry, nil
}

func (r *repository) FindByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) ([]*entities.SalaryEntry, error) {
	var entries []*entities.SalaryEntry

	// Order by: additions first, then discounts, then by created_at
	// Using CASE to put 'addition' before 'discount' in sorting
	if err := r.db.WithContext(ctx).
		Where("income_id = ? AND month = ? AND year = ?", incomeID, month, year).
		Order("CASE WHEN entry_type = 'addition' THEN 1 ELSE 2 END").
		Order("created_at ASC").
		Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar lançamentos de salário: %v", err)
	}

	return entries, nil
}

func (r *repository) FindByIncomeAndCategory(ctx context.Context, incomeID string, month, year int, category string) (*entities.SalaryEntry, error) {
	var entry entities.SalaryEntry

	if err := r.db.WithContext(ctx).
		Where("income_id = ? AND month = ? AND year = ? AND category = ?", incomeID, month, year, category).
		First(&entry).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar lançamento de salário por categoria: %v", err)
	}

	return &entry, nil
}

func (r *repository) DeleteByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) error {
	return r.db.WithContext(ctx).
		Where("income_id = ? AND month = ? AND year = ?", incomeID, month, year).
		Delete(&entities.SalaryEntry{}).Error
}

func (r *repository) GetMonthlySummary(ctx context.Context, incomeID string, month, year int) (*MonthlySummary, error) {
	var summary MonthlySummary

	// Query to calculate summary with JOIN to incomes table for base amount
	query := `
		SELECT
			i.id as income_id,
			? as month,
			? as year,
			i.amount as base_amount,
			COALESCE(SUM(CASE WHEN se.entry_type = 'addition' THEN se.amount ELSE 0 END), 0) as total_additions,
			COALESCE(SUM(CASE WHEN se.entry_type = 'discount' THEN se.amount ELSE 0 END), 0) as total_discounts,
			i.amount +
			COALESCE(SUM(CASE WHEN se.entry_type = 'addition' THEN se.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN se.entry_type = 'discount' THEN se.amount ELSE 0 END), 0) as net_amount
		FROM incomes i
		LEFT JOIN salary_entries se ON se.income_id = i.id
			AND se.month = ?
			AND se.year = ?
		WHERE i.id = ?
		GROUP BY i.id, i.amount
	`

	if err := r.db.WithContext(ctx).
		Raw(query, month, year, month, year, incomeID).
		Scan(&summary).Error; err != nil {
		return nil, fmt.Errorf("erro ao calcular resumo mensal: %v", err)
	}

	// Set the month and year explicitly (in case they weren't set by the query)
	summary.Month = month
	summary.Year = year

	return &summary, nil
}
