package gateways

import (
	"context"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/budget"
)

type BudgetGateway interface {
	Create(ctx context.Context, budget models.Budget) error
	Update(ctx context.Context, budget models.Budget) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Budget, error)
	List(ctx context.Context, status string, description string, page models.PageRequest) ([]models.Budget, int64, error)
	GetBudgetsWithoutMovement(ctx context.Context) ([]models.Budget, error)
}

type budgetGateway struct {
	repo budget.Repository
}

func NewBudgetGateway(repo budget.Repository) BudgetGateway {
	return &budgetGateway{repo: repo}
}

func (g *budgetGateway) Create(ctx context.Context, budget models.Budget) error {
	entity := mappers.ToBudgetEntity(budget)
	return g.repo.Create(ctx, entity)
}

func (g *budgetGateway) Update(ctx context.Context, budget models.Budget) error {
	entity := mappers.ToBudgetEntity(budget)
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
	return mappers.ToBudgetModel(entity), nil
}

func (g *budgetGateway) List(ctx context.Context, status string, description string, page models.PageRequest) ([]models.Budget, int64, error) {
	entities, count, err := g.repo.List(ctx, status, description, page)
	if err != nil {
		return nil, 0, err
	}

	budgets := make([]models.Budget, len(entities))
	for i, entity := range entities {
		budgets[i] = mappers.ToBudgetModel(&entity)
	}
	return budgets, count, nil
}

func (bg *budgetGateway) GetBudgetsWithoutMovement(ctx context.Context) (models []models.Budget, err error) {
	entites, err := bg.repo.GetBudgetsWithoutMovement(ctx)

	if err != nil {
		return models, err
	}

	for _, entity := range entites {
		models = append(models, mappers.ToBudgetModel(&entity))
	}

	return
}
