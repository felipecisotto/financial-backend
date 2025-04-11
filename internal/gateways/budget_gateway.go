package gateways

import (
	"context"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/budget"
)

type BudgetGateway interface {
	Create(ctx context.Context, budget models.Budget) error
	Update(ctx context.Context, budget models.Budget) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Budget, error)
	List(ctx context.Context, page models.PageRequest) ([]models.Budget, int64, error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Budget, error)
	GetSummary(ctx context.Context, month time.Month, year int) (models.BudgetSummary, error)
}

type budgetGateway struct {
	repo budget.Repository
}

func NewBudgetGateway(repo budget.Repository) BudgetGateway {
	return &budgetGateway{repo: repo}
}

func (g *budgetGateway) Create(ctx context.Context, budget models.Budget) error {
	entity := &entities.Budget{
		ID:          budget.ID(),
		Amount:      budget.Amount(),
		Description: budget.Description(),
		EndDate:     budget.EndDate(),
		CreatedAt:   budget.CreatedAt(),
		UpdatedAt:   budget.UpdatedAt(),
	}
	return g.repo.Create(ctx, entity)
}

func (g *budgetGateway) Update(ctx context.Context, budget models.Budget) error {
	entity := &entities.Budget{
		ID:          budget.ID(),
		Amount:      budget.Amount(),
		Description: budget.Description(),
		EndDate:     budget.EndDate(),
		CreatedAt:   budget.CreatedAt(),
		UpdatedAt:   time.Now(),
	}
	return g.repo.Update(ctx, entity)
}

func (g *budgetGateway) Delete(ctx context.Context, id string) error {
	return g.repo.Delete(ctx, id)
}

func (g *budgetGateway) Get(ctx context.Context, id string) (models.Budget, error) {
	entity, err := g.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return g.toModel(entity), nil
}

func (g *budgetGateway) List(ctx context.Context, page models.PageRequest) ([]models.Budget, int64, error) {
	entities, count, err := g.repo.List(ctx, page)
	if err != nil {
		return nil, 0, err
	}

	budgets := make([]models.Budget, len(entities))
	for i, entity := range entities {
		budgets[i] = g.toModel(&entity)
	}
	return budgets, count, nil
}

func (g *budgetGateway) ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Budget, error) {
	entities, err := g.repo.ListByMonth(ctx, month, year)
	if err != nil {
		return nil, err
	}

	budgets := make([]models.Budget, len(entities))
	for i, entity := range entities {
		budgets[i] = g.toModel(entity)
	}
	return budgets, nil
}

func (g *budgetGateway) GetSummary(ctx context.Context, month time.Month, year int) (models.BudgetSummary, error) {
	entities, err := g.repo.ListByMonth(ctx, month, year)
	if err != nil {
		return nil, err
	}

	var totalBudgeted float64
	var totalSpent float64

	for _, entity := range entities {
		totalBudgeted += entity.Amount
		for _, expense := range entity.Expenses {
			totalSpent += expense.Amount
		}
	}

	return models.NewBudgetSummary(totalBudgeted, totalSpent), nil
}

func (g *budgetGateway) toModel(entity *entities.Budget) models.Budget {
	return models.NewBudget(
		entity.ID,
		entity.Amount,
		entity.Description,
		entity.EndDate,
	)
}
