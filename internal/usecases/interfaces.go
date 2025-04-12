package usecases

import (
	"context"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
)

type ExpenseUseCase interface {
	Create(ctx context.Context, input *dtos.CreateExpenseRequest) (*dtos.ExpenseResponse, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*dtos.ExpenseResponse, error)
	List(ctx context.Context, input *dtos.ListExpensesRequest) (*dtos.ListExpensesResponse, error)
	ListByMonth(ctx context.Context, input *dtos.ListExpensesByMonthRequest) (*dtos.ListExpensesResponse, error)
	AssignToBudget(ctx context.Context, expenseID string, budgetID string) error
	RemoveFromBudget(ctx context.Context, expenseID string, budgetID string) error
}

type BudgetUseCase interface {
	Create(ctx context.Context, input *dtos.CreateBudgetRequest) (*dtos.BudgetResponse, error)
	Update(ctx context.Context, id string, input *dtos.UpdateBudgetRequest) (*dtos.BudgetResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*dtos.BudgetResponse, error)
	List(ctx context.Context, status string, description string, page models.PageRequest) (*models.Page[*dtos.BudgetResponse], error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]*dtos.BudgetResponse, error)
	GetSummary(ctx context.Context, month time.Month, year int) (*dtos.BudgetSummaryResponse, error)
}

type IncomeUseCase interface {
	Create(ctx context.Context, input *dtos.CreateIncomeRequest) (*dtos.IncomeResponse, error)
	Update(ctx context.Context, id string, input *dtos.UpdateIncomeRequest) (*dtos.IncomeResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*dtos.IncomeResponse, error)
	List(ctx context.Context) (*models.Page[*dtos.IncomeResponse], error)
	ListByType(ctx context.Context, incomeType string) ([]*dtos.IncomeResponse, error)
}
