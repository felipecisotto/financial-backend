package gateways

import (
	. "context"

	"financial-backend/internal/entities"
	. "financial-backend/internal/models"
	. "financial-backend/internal/repositories/income"
)

type IncomeGateway interface {
	Create(ctx Context, income Income) error
	Update(ctx Context, income Income) error
	Delete(ctx Context, id string) error
	Get(ctx Context, id string) (Income, error)
	List(ctx Context, incomeType, description string, page PageRequest) ([]Income, int64, error)
	SummaryByMonth(ctx Context, month, year int) (amount float64, err error)
}
type incomeGateway struct {
	repo Repository
}

func NewIncomeGateway(repo Repository) IncomeGateway {
	return &incomeGateway{repo: repo}
}

func (g *incomeGateway) Create(ctx Context, income Income) error {
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

func (g *incomeGateway) Update(ctx Context, income Income) error {
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

func (g *incomeGateway) Delete(ctx Context, id string) error {
	return g.repo.Delete(ctx, id)
}

func (g *incomeGateway) Get(ctx Context, id string) (Income, error) {
	entity, err := g.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return g.toModel(entity), nil
}

func (g *incomeGateway) List(ctx Context, incomeType, description string, page PageRequest) ([]Income, int64, error) {
	entities, count, err := g.repo.List(ctx, description, incomeType, int(page.Limit), page.Offset())
	if err != nil {
		return nil, 0, err
	}

	incomes := make([]Income, len(entities))
	for i, entity := range entities {
		incomes[i] = g.toModel(entity)
	}
	return incomes, count, nil
}

func (g *incomeGateway) toModel(entity *entities.Income) Income {
	income, _ := NewIncome(
		entity.ID,
		entity.Description,
		entity.Amount,
		IncomeType(entity.Type),
		entity.DueDay,
		entity.StartDate,
		entity.EndDate,
	)
	return income
}

func (g *incomeGateway) SummaryByMonth(ctx Context, month, year int) (amount float64, err error) {
	return g.repo.SummaryByMonth(ctx, month, year)
}
