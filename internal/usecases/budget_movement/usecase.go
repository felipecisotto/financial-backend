package budgetmovement

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/gateways"
	"financial-backend/internal/models"
)

type UseCase interface {
	FindByID(ctx context.Context, id string)
	Create(ctx context.Context, request dtos.BudgetMovementRequest) (dtos.BudgetMovementResponse, error)
	Find(ctx context.Context, params dtos.BudgetMovementParams) (models.Page[dtos.BudgetMovementResponse], error)
	CreateExpenseMovement(ctx context.Context, expense models.Expense) error
}

type useCase struct {
	gateway      gateways.BudgetMovementGateway
	budgetGatway gateways.BudgetGateway
}

func NewBudgetMovementUseCase(
	gateway gateways.BudgetMovementGateway,
	budgetGateway gateways.BudgetGateway,
) UseCase {
	return &useCase{
		budgetGatway: budgetGateway,
		gateway:      gateway,
	}
}
