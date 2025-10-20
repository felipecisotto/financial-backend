package dashboard

import (
	"context"
	"financial-backend/internal/dto/response"
	"financial-backend/internal/gateways"
	"fmt"
	"sync"
)

type UseCase interface {
	GetSummary(ctx context.Context, month, year int) (response.SummaryView, error)
	SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []response.SummaryBudgetUtilization, err error)
	GetMonthlyEvolution(ctx context.Context, year int) (response.MonthlyEvolutionView, error)
}

type useCase struct {
	expenseGateway        gateways.ExpenseGateway
	incomeGateway         gateways.IncomeGateway
	budgetMovementGateway gateways.BudgetMovementGateway
}

func (u useCase) GetSummary(ctx context.Context, month, year int) (response.SummaryView, error) {
	var wg sync.WaitGroup
	var income, expense float64
	var incomeErr, expenseErr error

	// Chama o SummaryByMonth para a renda em paralelo
	wg.Add(1)
	go func() {
		defer wg.Done()
		income, incomeErr = u.incomeGateway.SummaryByMonth(ctx, month, year)
	}()

	// Chama o SummaryByMonth para a despesa em paralelo
	wg.Add(1)
	go func() {
		defer wg.Done()
		expense, expenseErr = u.expenseGateway.SummaryByMonth(ctx, month, year)
	}()

	// Espera ambos os métodos terminarem
	wg.Wait()

	// Verifica se houve erro em qualquer uma das chamadas
	if incomeErr != nil {
		return response.SummaryView{}, incomeErr
	}
	if expenseErr != nil {
		return response.SummaryView{}, expenseErr
	}

	// Retorna o resumo com os valores de renda, despesa e restante
	return response.SummaryView{
		TotalIncome:    income,
		TotalExpense:   expense,
		TotalRemaining: income - expense,
	}, nil
}
func (u *useCase) SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []response.SummaryBudgetUtilization, err error) {
	data, err = u.budgetMovementGateway.SummaryBudgetUsageByMonthYear(ctx, month, year)
	return
}

func (u *useCase) GetMonthlyEvolution(ctx context.Context, year int) (response.MonthlyEvolutionView, error) {
	var wg sync.WaitGroup
	var incomes, expenses map[string]float64
	var incomeErr, expenseErr error

	// Fetch incomes in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		incomes, incomeErr = u.incomeGateway.MonthlyEvolution(ctx, year)
	}()

	// Fetch expenses in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		expenses, expenseErr = u.expenseGateway.MonthlyEvolution(ctx, year)
	}()

	// Wait for both operations to complete
	wg.Wait()

	// Check for errors
	if incomeErr != nil {
		return nil, incomeErr
	}
	if expenseErr != nil {
		return nil, expenseErr
	}

	// Build the result with all 12 months
	result := make(response.MonthlyEvolutionView, 0, 12)

	for month := 1; month <= 12; month++ {
		monthKey := fmt.Sprintf("%d", month)

		monthData := response.MonthlyData{
			Month:   month,
			Income:  incomes[monthKey],
			Expense: expenses[monthKey],
		}

		result = append(result, monthData)
	}

	return result, nil
}
func NewDashBoardUseCase(
	expenseGateway gateways.ExpenseGateway,
	incomeGateway gateways.IncomeGateway,
	budgetMovement gateways.BudgetMovementGateway,
) UseCase {
	return &useCase{
		expenseGateway:        expenseGateway,
		incomeGateway:         incomeGateway,
		budgetMovementGateway: budgetMovement,
	}
}
