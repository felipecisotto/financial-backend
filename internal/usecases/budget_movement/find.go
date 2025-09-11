package budgetmovement

import (
	"context"
	"financial-backend/internal/dto"
	"financial-backend/internal/mappers"
	"financial-backend/internal/domain/models"
	"math"
)

func (uc *useCase) FindByID(ctx context.Context, id string) {

}

func (uc *useCase) Find(ctx context.Context, params dto.BudgetMovementParams) (models.Page[dto.BudgetMovementResponse], error) {
	budgetMovements, count, err := uc.gateway.List(
		ctx,
		params.BudgetId,
		params.MovementType,
		params.Origin,
		params.Month,
		params.Year,
		models.PageRequest{
			Limit: params.Limit,
			Page:  params.Page,
		},
	)

	if err != nil {
		return models.Page[dto.BudgetMovementResponse]{}, err
	}

	budgetMovementResponses := make([]dto.BudgetMovementResponse, len(budgetMovements))

	for i, budgetMovement := range budgetMovements {
		budgetMovementResponses[i] = mappers.ToBudgetMovementDTO(budgetMovement)
	}

	return models.Page[dto.BudgetMovementResponse]{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int64(math.Ceil(float64(count) / float64(params.Limit))),
		Results:    budgetMovementResponses,
	}, nil
}
