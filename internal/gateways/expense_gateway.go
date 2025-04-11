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
	List(ctx context.Context) ([]models.Expense, error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Expense, error)
	AssignToBudget(ctx context.Context, expenseID string, budgetID string) error
	RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error
}

type expenseGateway struct {
	repo expense.Repository
}

func NewExpenseGateway(repo expense.Repository) ExpenseGateway {
	return &expenseGateway{repo: repo}
}

func (g *expenseGateway) Create(ctx context.Context, expense models.Expense) error {
	entity := &entities.Expense{
		ID:          expense.Id(),
		Description: expense.Description(),
		Amount:      expense.Amount(),
		Type:        string(expense.Type()),
		StartDate:   expense.StartDate(),
		DueDate:     expense.DueDate(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

func (g *expenseGateway) List(ctx context.Context) ([]models.Expense, error) {
	entities, err := g.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	expenses := make([]models.Expense, len(entities))
	for i, entity := range entities {
		expenses[i] = g.toModel(entity)
	}
	return expenses, nil
}

func (g *expenseGateway) ListByMonth(ctx context.Context, month time.Month, year int) ([]models.Expense, error) {
	entities, err := g.repo.ListByMonth(ctx, month, year)
	if err != nil {
		return nil, err
	}

	expenses := make([]models.Expense, len(entities))
	for i, entity := range entities {
		expenses[i] = g.toModel(entity)
	}
	return expenses, nil
}

func (g *expenseGateway) AssignToBudget(ctx context.Context, expenseID string, budgetID string) error {
	return g.repo.AssignToBudget(ctx, expenseID, budgetID)
}

func (g *expenseGateway) RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error {
	return g.repo.RemoveFromBudget(ctx, expenseID, budgetID)
}

func (g *expenseGateway) toModel(entity *entities.Expense) models.Expense {
	return models.NewExpense(entity.ID, entity.Description, entity.Amount, entity.Type, entity.DueDate, entity.StartDate)
}
