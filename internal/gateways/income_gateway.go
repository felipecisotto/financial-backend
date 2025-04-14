package gateways

import (
	"context"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/income"
)

type IncomeGateway interface {
	Create(ctx context.Context, income models.Income) error
	Update(ctx context.Context, income models.Income) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Income, error)
	List(ctx context.Context, incomeType, description string, page models.PageRequest) ([]models.Income, int64, error)
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
		Type:        string(income.Type()),
		DueDay:      income.DueDay(),
		StartDate:   income.StartDate(),
		EndDate:     income.EndDate(),
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
		Type:        string(income.Type()),
		DueDay:      income.DueDay(),
		StartDate:   income.StartDate(),
		EndDate:     income.EndDate(),
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

func (g *incomeGateway) List(ctx context.Context, incomeType, description string, page models.PageRequest) ([]models.Income, int64, error) {
	entities, count, err := g.repo.List(ctx, description, incomeType, int(page.Limit), page.Offset())
	if err != nil {
		return nil, 0, err
	}

	incomes := make([]models.Income, len(entities))
	for i, entity := range entities {
		incomes[i] = g.toModel(entity)
	}
	return incomes, count, nil
}

func (g *incomeGateway) toModel(entity *entities.Income) models.Income {
	income, _ := models.NewIncome(
		entity.ID,
		entity.Description,
		entity.Amount,
		models.IncomeType(entity.Type),
		entity.DueDay,
		entity.StartDate,
		entity.EndDate,
	)
	return income
}
