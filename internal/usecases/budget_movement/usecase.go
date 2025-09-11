package budgetmovement

import (
	"context"
	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto"
	"financial-backend/internal/gateways"
)

type UseCase interface {
	FindByID(ctx context.Context, id string)
	Create(ctx context.Context, request dto.BudgetMovementRequest) (dto.BudgetMovementResponse, error)
	Find(ctx context.Context, params dto.BudgetMovementParams) (models.Page[dto.BudgetMovementResponse], error)
	CreateExpenseMovement(ctx context.Context, expense models.Expense) error
	CreateRecurrencyMovements(ctx context.Context) error
}

type useCase struct {
	gateway        gateways.BudgetMovementGateway
	budgetGatway   gateways.BudgetGateway
	expenseGateway gateways.ExpenseGateway
}

func NewBudgetMovementUseCase(
	gateway gateways.BudgetMovementGateway,
	budgetGateway gateways.BudgetGateway,
	expenseGateway gateways.ExpenseGateway,
) UseCase {
	return &useCase{
		budgetGatway:   budgetGateway,
		gateway:        gateway,
		expenseGateway: expenseGateway,
	}
}
