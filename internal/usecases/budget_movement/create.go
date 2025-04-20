package budgetmovement

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"
	"time"

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

func (uc *useCase) CreateRecurrencyMovements(ctx context.Context) error {
	var movements []models.BudgetMovement

	expenseMovements, err := uc.createExpenseRecurrencyMovements(ctx)

	if err != nil {
		return err
	}

	movements = append(movements, expenseMovements...)

	budgetStartMovements, err := uc.createBudgetStartMovements(ctx)

	if err != nil {
		return err
	}

	movements = append(movements, budgetStartMovements...)

	return uc.gateway.CreateAll(ctx, movements)
}

func (uc *useCase) createBudgetStartMovements(ctx context.Context) (movements []models.BudgetMovement, err error) {
	budgets, err := uc.budgetGatway.GetBudgetsWithoutMovement(ctx)
	if err != nil {
		return movements, err
	}

	for _, budget := range budgets {
		movements = append(movements, buildMovementByBudget(budget))
	}

	return
}

func (uc *useCase) createExpenseRecurrencyMovements(ctx context.Context) ([]models.BudgetMovement, error) {
	movements := []models.BudgetMovement{}
	expenses, err := uc.expenseGateway.GetExpensesWithoutMovimentInMonth(ctx)

	if err != nil {
		return make([]models.BudgetMovement, 0), err
	}

	actualDate := time.Now()

	for _, expense := range expenses {
		recurrency := expense.Recurrency()
		if *recurrency == models.ExpenseRecurrencyWeekly {
			for range uc.countWeekdayInMonth(actualDate.Year(), actualDate.Month(), time.Weekday(expense.DueDay())) {
				movements = append(movements, buildMovementByExpense(expense, int(actualDate.Month()), actualDate.Year(), *expense.Budget()))
			}
			continue
		}
		movements = append(movements, buildMovementByExpense(expense, int(actualDate.Month()), actualDate.Year(), *expense.Budget()))
	}

	return movements, nil
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

func buildMovementByBudget(budget models.Budget) models.BudgetMovement {
	return models.NewBudgetMovement(
		uuid.New().String(),
		budget.ID(),
		budget,
		budget.ID(),
		nil,
		int(time.Now().Month()),
		time.Now().Year(),
		models.MovementStart,
		int(budget.Amount()),
	)
}

func (uc *useCase) countWeekdayInMonth(year int, month time.Month, weekday time.Weekday) int {
	loc := time.UTC
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	lastDay := firstDay.AddDate(0, 1, -1)

	count := 0
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == weekday {
			count++
		}
	}

	return count
}
