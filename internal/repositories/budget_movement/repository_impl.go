package budgetmovement

import (
	"context"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(ctx context.Context, budgetMovement entities.BudgetMovement) error {
	return r.db.WithContext(ctx).Create(budgetMovement).Error
}

func (r *repository) CreateAll(ctx context.Context, budgetMovements []entities.BudgetMovement) error {
	return r.db.WithContext(ctx).CreateInBatches(budgetMovements, 50).Error
}

// GetById implements Repository.
func (r *repository) GetById(ctx context.Context, id string) (budget *entities.BudgetMovement, err error) {
	if err := r.db.WithContext(ctx).First(&budget, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return
}

// List implements Repository.
func (r *repository) List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) (budgets []entities.BudgetMovement, count int64, err error) {
	selectColumns := `SELECT 
		bm.id,
		bm.budget_id,
		bm.type,
		bm.origin,
		bm.month,
		bm.year,
		bm.amount,
		bm.created_at,
		COALESCE(i.description, e.description, b1.description) AS origin_description,
		b.id AS "budget__id",
	b.description AS "budget__description",
	b.amount AS "budget__amount",
	b.end_date AS "budget__end_date",
	b.created_at AS "budget__created_at",
	b.updated_at AS "budget__updated_at"`
	countColumns := `SELECT count(1)`
	query := `
	FROM budget_movements bm
	JOIN budgets b ON bm.budget_id = b.id
	LEFT JOIN incomes i ON bm.origin = i.id AND bm.type = 'income'
	LEFT JOIN expenses e ON bm.origin = e.id AND bm.type = 'expense'
	LEFT JOIN budgets b1 ON bm.origin = b1.id AND bm.type = 'budget'
	WHERE 1=1
	`

	var args []interface{}

	if budgetId != "" {
		query += " AND bm.budget_id = ?"
		args = append(args, budgetId)
	}

	if movementType != "" {
		query += " AND bm.type = ?"
		args = append(args, movementType)
	}

	if origin != "" {
		query += " AND bm.origin = ?"
		args = append(args, origin)
	}

	if month != 0 {
		query += " AND bm.month = ?"
		args = append(args, month)
	}

	if year != 0 {
		query += " AND bm.year = ?"
		args = append(args, year)
	}

	pagedQuery := query + " ORDER BY bm.created_at DESC LIMIT ? OFFSET ?"
	pagedArgs := append(args, page.Limit, page.Offset())

	if err := r.db.WithContext(ctx).Raw(selectColumns+pagedQuery, pagedArgs...).Preload("Budget").Find(&budgets).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar movimentações: %w", err)
	}

	if err := r.db.WithContext(ctx).Raw(countColumns+query, args...).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar movimentações: %w", err)
	}

	return
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
