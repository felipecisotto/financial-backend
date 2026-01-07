package gateways

import (
	"context"

	"financial-backend/internal/domain/entities"
	"financial-backend/internal/domain/models"
	"financial-backend/internal/mappers"
	"financial-backend/internal/repositories/salary_entry"
)

// SalaryEntryGateway defines the interface for salary entry gateway operations
type SalaryEntryGateway interface {
	Create(ctx context.Context, entry models.SalaryEntry) error
	CreateBatch(ctx context.Context, entries []models.SalaryEntry) error
	Update(ctx context.Context, entry models.SalaryEntry) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (models.SalaryEntry, error)
	FindByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) ([]models.SalaryEntry, error)
	FindByIncomeAndCategory(ctx context.Context, incomeID string, month, year int, category string) (models.SalaryEntry, error)
	DeleteByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) error
	GetMonthlySummary(ctx context.Context, incomeID string, month, year int) (*salary_entry.MonthlySummary, error)
}

type salaryEntryGateway struct {
	repo salary_entry.Repository
}

// NewSalaryEntryGateway creates a new instance of the salary entry gateway
func NewSalaryEntryGateway(repo salary_entry.Repository) SalaryEntryGateway {
	return &salaryEntryGateway{repo: repo}
}

func (g *salaryEntryGateway) Create(ctx context.Context, entry models.SalaryEntry) error {
	entity := mappers.SalaryEntryModelToEntity(entry)
	return g.repo.Create(ctx, entity)
}

func (g *salaryEntryGateway) CreateBatch(ctx context.Context, entries []models.SalaryEntry) error {
	entities := make([]*entities.SalaryEntry, len(entries))
	for i, entry := range entries {
		entities[i] = mappers.SalaryEntryModelToEntity(entry)
	}
	return g.repo.CreateBatch(ctx, entities)
}

func (g *salaryEntryGateway) Update(ctx context.Context, entry models.SalaryEntry) error {
	entity := mappers.SalaryEntryModelToEntity(entry)
	return g.repo.Update(ctx, entity)
}

func (g *salaryEntryGateway) Delete(ctx context.Context, id string) error {
	return g.repo.Delete(ctx, id)
}

func (g *salaryEntryGateway) FindByID(ctx context.Context, id string) (models.SalaryEntry, error) {
	entity, err := g.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.SalaryEntryEntityToModel(entity)
}

func (g *salaryEntryGateway) FindByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) ([]models.SalaryEntry, error) {
	entities, err := g.repo.FindByIncomeAndPeriod(ctx, incomeID, month, year)
	if err != nil {
		return nil, err
	}

	salaryEntries := make([]models.SalaryEntry, 0, len(entities))
	for _, entity := range entities {
		model, err := mappers.SalaryEntryEntityToModel(entity)
		if err != nil {
			return nil, err
		}
		salaryEntries = append(salaryEntries, model)
	}

	return salaryEntries, nil
}

func (g *salaryEntryGateway) FindByIncomeAndCategory(ctx context.Context, incomeID string, month, year int, category string) (models.SalaryEntry, error) {
	entity, err := g.repo.FindByIncomeAndCategory(ctx, incomeID, month, year, category)
	if err != nil {
		return nil, err
	}

	return mappers.SalaryEntryEntityToModel(entity)
}

func (g *salaryEntryGateway) DeleteByIncomeAndPeriod(ctx context.Context, incomeID string, month, year int) error {
	return g.repo.DeleteByIncomeAndPeriod(ctx, incomeID, month, year)
}

func (g *salaryEntryGateway) GetMonthlySummary(ctx context.Context, incomeID string, month, year int) (*salary_entry.MonthlySummary, error) {
	return g.repo.GetMonthlySummary(ctx, incomeID, month, year)
}
