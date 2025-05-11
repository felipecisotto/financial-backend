package dashboard

import (
	. "context"
	. "financial-backend/internal/gateways"
	"financial-backend/internal/views"
	"golang.org/x/net/context"
	"sync"
)

type UseCase interface {
	GetSummary(ctx Context, month, year int) (views.SummaryView, error)
	SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []views.SummaryBudgetUtilization, err error)
}

type useCase struct {
	expenseGateway        ExpenseGateway
	incomeGateway         IncomeGateway
	budgetMovementGateway BudgetMovementGateway
}

func (u useCase) GetSummary(ctx Context, month, year int) (views.SummaryView, error) {
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
		return views.SummaryView{}, incomeErr
	}
	if expenseErr != nil {
		return views.SummaryView{}, expenseErr
	}

	// Retorna o resumo com os valores de renda, despesa e restante
	return views.SummaryView{
		TotalIncome:    income,
		TotalExpense:   expense,
		TotalRemaining: income - expense,
	}, nil
}
func (u *useCase) SummaryBudgetUsageByMonthYear(ctx context.Context, month, year int) (data []views.SummaryBudgetUtilization, err error) {
	data, err = u.budgetMovementGateway.SummaryBudgetUsageByMonthYear(ctx, month, year)
	return
}
func NewDashBoardUseCase(
	expenseGateway ExpenseGateway,
	incomeGateway IncomeGateway,
	budgetMovement BudgetMovementGateway,
) UseCase {
	return &useCase{
		expenseGateway:        expenseGateway,
		incomeGateway:         incomeGateway,
		budgetMovementGateway: budgetMovement,
	}
}
