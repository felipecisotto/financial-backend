package gateways

import (
	"context"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/income"
)

type IncomeGateway interface {
	Create(ctx context.Context, income models.Income) error
	Update(ctx context.Context, income models.Income) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Income, error)
	List(ctx context.Context) ([]models.Income, error)
	ListByType(ctx context.Context, incomeType models.IncomeType) ([]models.Income, error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Income, error)
}

type incomeGateway struct {
	repo income.Repository
}

func NewIncomeGateway(repo income.Repository) IncomeGateway {
	return &incomeGateway{repo: repo}
}

func (g *incomeGateway) Create(ctx context.Context, income models.Income) error {
	entity := &entities.Income{
		ID:          income.ID(),
		Description: income.Description(),
		Amount:      income.Amount(),
		Type:        entities.IncomeType(income.Type()),
		DueDay:      income.DueDay(),
		CreatedAt:   income.CreatedAt(),
		UpdatedAt:   income.UpdatedAt(),
	}
	return g.repo.Create(ctx, entity)
}

func (g *incomeGateway) Update(ctx context.Context, income models.Income) error {
	entity := &entities.Income{
		ID:          income.ID(),
		Description: income.Description(),
		Amount:      income.Amount(),
		Type:        entities.IncomeType(income.Type()),
		DueDay:      income.DueDay(),
		CreatedAt:   income.CreatedAt(),
		UpdatedAt:   income.UpdatedAt(),
	}
	return g.repo.Update(ctx, entity)
}

func (g *incomeGateway) Delete(ctx context.Context, id string) error {
	return g.repo.Delete(ctx, id)
}

func (g *incomeGateway) Get(ctx context.Context, id string) (models.Income, error) {
	entity, err := g.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return g.toModel(entity), nil
}

func (g *incomeGateway) List(ctx context.Context) ([]models.Income, error) {
	entities, err := g.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	incomes := make([]models.Income, len(entities))
	for i, entity := range entities {
		incomes[i] = g.toModel(entity)
	}
	return incomes, nil
}

func (g *incomeGateway) ListByType(ctx context.Context, incomeType models.IncomeType) ([]models.Income, error) {
	entities, err := g.repo.ListByType(ctx, entities.IncomeType(incomeType))
	if err != nil {
		return nil, err
	}

	incomes := make([]models.Income, len(entities))
	for i, entity := range entities {
		incomes[i] = g.toModel(entity)
	}
	return incomes, nil
}

func (g *incomeGateway) ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Income, error) {
	entities, err := g.repo.ListByMonth(ctx, month, year)
	if err != nil {
		return nil, err
	}

	incomes := make([]models.Income, len(entities))
	for i, entity := range entities {
		incomes[i] = g.toModel(entity)
	}
	return incomes, nil
}

func (g *incomeGateway) toModel(entity *entities.Income) models.Income {
	return models.NewIncome(
		entity.ID,
		entity.Description,
		entity.Amount,
		models.IncomeType(entity.Type),
		entity.DueDay,
	)
}
