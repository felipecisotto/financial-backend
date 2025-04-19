package budgetmovement

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"

	"github.com/google/uuid"
)

func (uc *useCase) Create(ctx context.Context, request dtos.BudgetMovementRequest) (dtos.BudgetMovementResponse, error) {
	budgetMovement := mappers.FromDTOToBudgetMovementModel(request)
	err := uc.gateway.Create(ctx, budgetMovement)

	if err != nil {
		return dtos.BudgetMovementResponse{}, err
	}

	return mappers.ToBudgetMovementDTO(budgetMovement), nil
}

func (uc *useCase) CreateExpenseMovement(ctx context.Context, expense models.Expense) error {
	if expense.BudgetId() == nil {
		return nil
	}
	var budget models.Budget

	if expense.Budget() == nil {
		foundBudget, err := uc.budgetGatway.Get(ctx, *expense.BudgetId())
		if err != nil {
			return err
		}
		budget = foundBudget
	} else {
		budget = *expense.Budget()
	}

	if expense.Installments() == nil {
		return uc.gateway.Create(ctx, buildMovementByExpense(expense, int(expense.StartDate().Month()), expense.StartDate().Year(), budget))
	}

	for i := 0; i < *expense.Installments(); i++ {
		date := expense.StartDate().AddDate(0, i, 0)
		if err := uc.gateway.Create(ctx, buildMovementByExpense(expense, int(date.Month()), date.Year(), budget)); err != nil {
			return err
		}
	}
	return nil
}

func buildMovementByExpense(expense models.Expense, month, year int, budget models.Budget) models.BudgetMovement {
	return models.NewBudgetMovement(
		uuid.New().String(),
		*expense.BudgetId(),
		budget,
		expense.Id(),
		nil,
		month,
		year,
		models.MovementExpense,
		int(expense.Amount()),
	)
}
