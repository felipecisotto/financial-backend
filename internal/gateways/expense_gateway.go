package gateways

import (
	"context"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/mappers"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/expense"
)

type ExpenseGateway interface {
	Create(ctx context.Context, expense models.Expense) error
	Update(ctx context.Context, expense models.Expense) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Expense, error)
	List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string, page models.PageRequest) ([]models.Expense, int64, error)
	GetExpensesWithoutMovementInMonth(ctx context.Context) ([]models.Expense, error)
	SummaryByMonth(ctx context.Context, month, year int) (amount float64, err error)
}

type expenseGateway struct {
	repo expense.Repository
}

func NewExpenseGateway(repo expense.Repository) ExpenseGateway {
	return &expenseGateway{repo: repo}
}

func (g *expenseGateway) Create(ctx context.Context, expense models.Expense) error {
	entity := &entities.Expense{
		ID:           expense.Id(),
		Description:  expense.Description(),
		Amount:       expense.Amount(),
		Type:         string(expense.Type()),
		BudgetID:     expense.BudgetId(),
		Recurrency:   (*string)(expense.Recurrency()),
		Method:       string(expense.Method()),
		Installments: expense.Installments(),
		StartDate:    expense.StartDate(),
		DueDay:       expense.DueDay(),
		EndDate:      expense.EndDate(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return g.repo.Create(ctx, entity)
}

func (g *expenseGateway) Update(ctx context.Context, expense models.Expense) error {
	entity := &entities.Expense{
		ID:          expense.Id(),
		Description: expense.Description(),
		Amount:      expense.Amount(),
		Type:        string(expense.Type()),
		UpdatedAt:   time.Now(),
	}
	return g.repo.Update(ctx, entity)
}

func (g *expenseGateway) Delete(ctx context.Context, id string) error {
	return g.repo.Delete(ctx, id)
}

func (g *expenseGateway) Get(ctx context.Context, id string) (models.Expense, error) {
	entity, err := g.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToExpenseModel(entity), nil
}

func (g *expenseGateway) List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string, page models.PageRequest) ([]models.Expense, int64, error) {
	entities, count, err := g.repo.List(ctx, description, expenseType, category, budgetId, recurrecy, method, page)
	if err != nil {
		return nil, 0, err
	}

	expenses := make([]models.Expense, len(entities))
	for i, entity := range entities {
		expenses[i] = mappers.ToExpenseModel(entity)
	}
	return expenses, count, nil
}

func (g *expenseGateway) GetExpensesWithoutMovementInMonth(ctx context.Context) ([]models.Expense, error) {
	entities, err := g.repo.GetExpensesWithoutMovimentInMonth(ctx)
	responses := make([]models.Expense, len(entities))

	if err != nil {
		return responses, err
	}

	for i, entity := range entities {
		responses[i] = mappers.ToExpenseModel(entity)
	}

	return responses, nil
}
func (g *expenseGateway) SummaryByMonth(ctx context.Context, month, year int) (amount float64, err error) {
	amount, err = g.repo.SummaryByMonth(ctx, month, year)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
