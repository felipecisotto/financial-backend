package budgetmovement

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/mappers"
)

func (uc *useCase) Create(ctx context.Context, request dtos.BudgetMovementRequest) (dtos.BudgetMovementResponse, error) {
	budgetMovement := mappers.FromDTOToBudgetMovementModel(request)
	err := uc.gateway.Create(ctx, budgetMovement)

	if err != nil {
		return dtos.BudgetMovementResponse{}, err
	}

	return mappers.ToBudgetMovementDTO(budgetMovement), nil
}
