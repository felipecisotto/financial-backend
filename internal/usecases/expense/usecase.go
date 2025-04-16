package expense

import (
	"context"
	"fmt"
	"math"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/gateways"
	"financial-backend/internal/models"

	"github.com/google/uuid"
)

type useCase struct {
	expenseGateway gateways.ExpenseGateway
	budgetGateway  gateways.BudgetGateway
	defaultDueDate int
}

type UseCase interface {
	Create(ctx context.Context, input *dtos.ExpenseDTO) (*dtos.ExpenseResponse, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*dtos.ExpenseResponse, error)
	List(ctx context.Context, input *dtos.ListExpensesRequest) (*models.Page[*dtos.ExpenseResponse], error)
}

func NewUseCase(expenseGateway gateways.ExpenseGateway, budgetGateway gateways.BudgetGateway, defaultDueDate int) UseCase {
	return &useCase{
		expenseGateway: expenseGateway,
		budgetGateway:  budgetGateway,
		defaultDueDate: defaultDueDate,
	}
}

func (uc *useCase) Create(ctx context.Context, input *dtos.ExpenseDTO) (*dtos.ExpenseResponse, error) {
	var newEndDate *time.Time

	if input.Installments != nil && input.EndDate == nil {
		newEndDateValue := time.Date(time.Now().Year(), time.Now().Month()+time.Month(*input.Installments), uc.defaultDueDate, 0, 0, 0, 0, time.UTC)
		newEndDate = &newEndDateValue // Take the address of the new time.Time value
	} else {
		newEndDate = input.EndDate
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
		input.StartDate,
		newEndDate,
		nil,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.expenseGateway.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("erro ao criar despesa: %v", err)
	}

	return uc.toExpenseResponse(expense), nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.expenseGateway.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao excluir despesa: %v", err)
	}
	return nil
}

func (uc *useCase) FindByID(ctx context.Context, id string) (*dtos.ExpenseResponse, error) {
	expense, err := uc.expenseGateway.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar despesa: %v", err)
	}
	return uc.toExpenseResponse(expense), nil
}

func (uc *useCase) List(ctx context.Context, request *dtos.ListExpensesRequest) (*models.Page[*dtos.ExpenseResponse], error) {
	expenses, count, err := uc.expenseGateway.List(
		ctx,
		request.Description,
		request.Type,
		request.Category,
		request.BudgetID,
		request.Recurrency,
		request.Method,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas: %v", err)
	}

	responses := make([]*dtos.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		responses[i] = uc.toExpenseResponse(expense)
	}

	return &models.Page[*dtos.ExpenseResponse]{
		Page:       request.Page,
		Limit:      request.Limit,
		TotalPages: int64(math.Ceil(float64(count) / float64(request.Limit))),
		Results:    responses,
	}, nil
}

func (uc *useCase) toExpenseResponse(expense models.Expense) *dtos.ExpenseResponse {
	var budget *dtos.BudgetResponse
	if expense.Budget() != nil {
		budgetModel := expense.Budget()
		budget = &dtos.BudgetResponse{
			ID:          (*budgetModel).ID(),
			Amount:      (*budgetModel).Amount(),
			Description: (*budgetModel).Description(),
			EndDate:     (*budgetModel).EndDate(),
			Status:      string((*budgetModel).Status()),
			CreatedAt:   (*budgetModel).CreatedAt(),
			UpdatedAt:   (*budgetModel).UpdatedAt(),
		}
	}

	return &dtos.ExpenseResponse{
		ID: expense.Id(),
		ExpenseDTO: dtos.ExpenseDTO{
			Description:  expense.Description(),
			Amount:       expense.Amount(),
			Type:         string(expense.Type()),
			BudgetID:     expense.BudgetId(),
			Recurrency:   (*string)(expense.Recurrency()),
			Method:       string(expense.Method()),
			Installments: expense.Installments(),
			DueDay:       expense.DueDay(),
			StartDate:    expense.StartDate(),
			EndDate:      expense.EndDate(),
			Budget:       budget,
		},
	}
}
