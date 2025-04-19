package expense

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"financial-backend/internal/models/events"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (uc *useCase) Create(ctx context.Context, input *dtos.ExpenseDTO) (*dtos.ExpenseResponse, error) {
	var endDate *time.Time
	var startDate time.Time

	startDate = input.StartDate
	endDate = input.EndDate

	if input.Method == string(models.ExpenseMethodCreditCard) {
		if startDate.Day() > uc.defaultDueDate {
			startDate = startDate.AddDate(0, 1, -startDate.Day()+uc.defaultDueDate)
		} else {
			startDate = startDate.AddDate(0, 0, uc.defaultDueDate-startDate.Day())
		}

		endDate = new(time.Time)
		*endDate = startDate.AddDate(0, *input.Installments-1, 0)
	}

	expense, err := models.NewExpense(
		uuid.New().String(),
		input.Description,
		input.Amount,
		input.Type,
		input.BudgetID,
		input.Recurrency,
		input.Method,
		input.Installments,
		input.DueDay,
		startDate,
		endDate,
		nil,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.expenseGateway.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("erro ao criar despesa: %v", err)
	}

	uc.eventPublisher.Publish(&events.ExpenseCreatedEvent{
		Expense: expense,
		Context: ctx,
	})

	return uc.toExpenseResponse(expense), nil
}
