package gateways

import (
	"context"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/expense"
)

type ExpenseGateway interface {
	Create(ctx context.Context, expense models.Expense) error
	Update(ctx context.Context, expense models.Expense) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.Expense, error)
	List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string) ([]models.Expense, int64, error)
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
	return g.toModel(entity), nil
}

func (g *expenseGateway) List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string) ([]models.Expense, int64, error) {
	entities, count, err := g.repo.List(ctx, description, expenseType, category, budgetId, recurrecy, method)
	if err != nil {
		return nil, 0, err
	}

	expenses := make([]models.Expense, len(entities))
	for i, entity := range entities {
		expenses[i] = g.toModel(entity)
	}
	return expenses, count, nil
}

func (g *expenseGateway) toModel(entity *entities.Expense) models.Expense {
	var budget *models.Budget
	if entity.Budget != nil {
		newBudget := models.NewBudget( // Create a new instance
			entity.ID,
			entity.Amount,
			entity.Description,
			entity.EndDate,
		)
		budget = &newBudget
	}
	expense, _ := models.NewExpense(
		entity.ID,
		entity.Description,
		entity.Amount,
		entity.Type,
		entity.BudgetID,
		entity.Recurrency,
		entity.Method,
		entity.Installments,
		entity.DueDay,
		entity.StartDate,
		entity.EndDate,
		budget,
	)
	return expense
}
