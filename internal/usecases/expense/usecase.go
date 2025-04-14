package expense

import (
	"context"
	"fmt"
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
	List(ctx context.Context, input *dtos.ListExpensesRequest) (*dtos.ListExpensesResponse, error)
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
	}

	newEndDate = input.EndDate

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

func (uc *useCase) List(ctx context.Context, request *dtos.ListExpensesRequest) (*dtos.ListExpensesResponse, error) {
	expenses, err := uc.expenseGateway.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar despesas: %v", err)
	}

	var response []dtos.ExpenseResponse
	for _, expense := range expenses {
		response = append(response, *uc.toExpenseResponse(expense))
	}

	return &dtos.ListExpensesResponse{
		Expenses: response,
		Total:    int64(len(response)),
	}, nil
}

func (uc *useCase) toExpenseResponse(expense models.Expense) *dtos.ExpenseResponse {
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
		},
	}
}
