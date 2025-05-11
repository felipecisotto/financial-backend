package gateways

import (
	"context"
	"financial-backend/internal/entities"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"
	budgetmovementRepository "financial-backend/internal/repositories/budget_movement"
	. "financial-backend/internal/views"
)

type BudgetMovementGateway interface {
	Create(ctx context.Context, budgetMovement models.BudgetMovement) error
	CreateAll(ctx context.Context, movements []models.BudgetMovement) error
	List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) ([]models.BudgetMovement, int64, error)
	GetByID(ctx context.Context, id string) (models.BudgetMovement, error)
	SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []SummaryBudgetUtilization, err error)
}

type budgetMovementGateway struct {
	repository budgetmovementRepository.Repository
}

func NewBudgetMovementGateway(repository budgetmovementRepository.Repository) BudgetMovementGateway {
	return &budgetMovementGateway{
		repository: repository,
	}
}

// Create implements BudgetMovementGateway.
func (b *budgetMovementGateway) Create(ctx context.Context, budgetMovement models.BudgetMovement) error {
	return b.repository.Create(ctx, mappers.ToBudgetMovementEntity(budgetMovement))
}

// GetByID implements BudgetMovementGateway.
func (b *budgetMovementGateway) GetByID(ctx context.Context, id string) (models.BudgetMovement, error) {
	entity, err := b.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ToBudgetMovementModel(*entity), nil
}

// List implements BudgetMovementGateway.
func (b *budgetMovementGateway) List(ctx context.Context, budgetId, movementType, origin string, month, year int, page models.PageRequest) ([]models.BudgetMovement, int64, error) {
	entities, count, err := b.repository.List(ctx, budgetId, movementType, origin, month, year, page)

	if err != nil {
		return nil, 0, err
	}
	responses := make([]models.BudgetMovement, len(entities))

	for i, entity := range entities {
		responses[i] = mappers.ToBudgetMovementModel(entity)
	}
	return responses, count, err
}

func (b *budgetMovementGateway) CreateAll(ctx context.Context, movements []models.BudgetMovement) error {
	entities := make([]entities.BudgetMovement, len(movements))

	for i, model := range movements {
		entities[i] = mappers.ToBudgetMovementEntity(model)
	}

	return b.repository.CreateAll(ctx, entities)
}

func (b *budgetMovementGateway) SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []SummaryBudgetUtilization, err error) {
	data, err = b.repository.SummaryBudgetUsageByMonthYear(ctx, month, year)
	if err != nil {
		return []SummaryBudgetUtilization{}, err
	}

	return data, nil
}
