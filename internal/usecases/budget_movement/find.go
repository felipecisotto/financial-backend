package budgetmovement

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"
	"math"
)

func (uc *useCase) FindByID(ctx context.Context, id string) {

}

func (uc *useCase) Find(ctx context.Context, params dtos.BudgetMovementParams) (models.Page[dtos.BudgetMovementResponse], error) {
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
		return models.Page[dtos.BudgetMovementResponse]{}, err
	}

	budgetMovementResponses := make([]dtos.BudgetMovementResponse, len(budgetMovements))

	for i, budgetMovement := range budgetMovements {
		budgetMovementResponses[i] = mappers.ToBudgetMovementDTO(budgetMovement)
	}

	return models.Page[dtos.BudgetMovementResponse]{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int64(math.Ceil(float64(count) / float64(params.Limit))),
		Results:    budgetMovementResponses,
	}, nil
}
