package mappers

import (
	"financial-backend/internal/dtos"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"

	"github.com/google/uuid"
)

func ToBudgetModel(entity *entities.Budget) models.Budget {
	return models.NewBudget(
		entity.ID,
		entity.Amount,
		entity.Description,
		entity.EndDate,
	)
}

func ToBudgetEntity(budget models.Budget) *entities.Budget {
	return &entities.Budget{
		ID:          budget.ID(),
		Amount:      budget.Amount(),
		Description: budget.Description(),
		EndDate:     budget.EndDate(),
		CreatedAt:   budget.CreatedAt(),
		UpdatedAt:   budget.UpdatedAt(),
	}
}

func FromDTOToBudgetModel(dto dtos.CreateBudgetRequest) models.Budget {
	return models.NewBudget(
		uuid.New().String(),
		dto.Amount,
		dto.Description,
		dto.EndDate,
	)
}

func ToBudgetResponse(budget models.Budget) dtos.BudgetResponse {
	return dtos.BudgetResponse{
		ID:          budget.ID(),
		Amount:      budget.Amount(),
		Description: budget.Description(),
		EndDate:     budget.EndDate(),
		Status:      string(budget.Status()),
		CreatedAt:   budget.CreatedAt(),
		UpdatedAt:   budget.UpdatedAt(),
	}
}
