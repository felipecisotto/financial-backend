package expense

import (
	"context"
	"fmt"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/gateways"
	"financial-backend/internal/models"
	"financial-backend/internal/usecases"

	"github.com/google/uuid"
)

type useCase struct {
	expenseGateway gateways.ExpenseGateway
	budgetGateway  gateways.BudgetGateway
	defaultDueDate int
}

func NewUseCase(expenseGateway gateways.ExpenseGateway, budgetGateway gateways.BudgetGateway, defaultDueDate int) usecases.ExpenseUseCase {
	return &useCase{
		expenseGateway: expenseGateway,
		budgetGateway:  budgetGateway,
		defaultDueDate: defaultDueDate,
	}
}

func (uc *useCase) Create(ctx context.Context, input *dtos.CreateExpenseRequest) (*dtos.ExpenseResponse, error) {

	dueDate := time.Date(time.Now().Year(), time.Now().Month()+time.Month(*input.Installments), uc.defaultDueDate, 0, 0, 0, 0, time.UTC)
	expense := models.NewExpense(uuid.New().String(), input.Description, input.Amount, input.Type, time.Now(), dueDate)

	if err := uc.expenseGateway.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("erro ao criar despesa: %v", err)
	}

	// Se um orçamento foi especificado, vincula a despesa a ele
	if input.BudgetID != nil && *input.BudgetID != "" {
		if err := uc.AssignToBudget(ctx, expense.Id(), *input.BudgetID); err != nil {
			return nil, fmt.Errorf("erro ao vincular despesa ao orçamento: %v", err)
		}
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

func (uc *useCase) ListByMonth(ctx context.Context, request *dtos.ListExpensesByMonthRequest) (*dtos.ListExpensesResponse, error) {
	month := request.Month.Month()
	year := request.Month.Year()

	expenses, err := uc.expenseGateway.ListByMonth(ctx, month, year)
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

func (uc *useCase) AssignToBudget(ctx context.Context, expenseID string, budgetID string) error {
	// Verifica se o orçamento existe
	_, err := uc.budgetGateway.Get(ctx, budgetID)
	if err != nil {
		return fmt.Errorf("erro ao buscar orçamento: %v", err)
	}

	if err := uc.expenseGateway.AssignToBudget(ctx, expenseID, budgetID); err != nil {
		return fmt.Errorf("erro ao vincular despesa ao orçamento: %v", err)
	}

	return nil
}

func (uc *useCase) RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error {
	if err := uc.expenseGateway.RemoveFromBudget(ctx, expenseID, budgetID); err != nil {
		return fmt.Errorf("erro ao desvincular despesa do orçamento: %v", err)
	}
	return nil
}

func (uc *useCase) toExpenseResponse(expense models.Expense) *dtos.ExpenseResponse {
	startDate := expense.StartDate()
	dueDate := expense.DueDate()

	return &dtos.ExpenseResponse{
		ID:          expense.Id(),
		Description: expense.Description(),
		Amount:      expense.Amount(),
		Type:        string(expense.Type()),
		StartDate:   &startDate, 
		DueDate:     &dueDate,
	}
}

