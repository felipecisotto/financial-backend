package budget

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/gateways"
	"financial-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, dto *dtos.CreateBudgetRequest) (*dtos.BudgetResponse, error)
	Update(ctx context.Context, id string, dto *dtos.UpdateBudgetRequest) (*dtos.BudgetResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*dtos.BudgetResponse, error)
	List(ctx context.Context, page models.PageRequest) (*models.Page[*dtos.BudgetResponse], error)
	ListByMonth(ctx context.Context, month time.Month, year int) ([]*dtos.BudgetResponse, error)
	GetSummary(ctx context.Context, month time.Month, year int) (*dtos.BudgetSummaryResponse, error)
}

type useCase struct {
	gateway gateways.BudgetGateway
}

func NewUseCase(gateway gateways.BudgetGateway) UseCase {
	return &useCase{gateway: gateway}
}

func (uc *useCase) Create(ctx context.Context, dto *dtos.CreateBudgetRequest) (*dtos.BudgetResponse, error) {
	budget := models.NewBudget(
		uuid.New().String(),
		dto.Amount,
		dto.Description,
		dto.EndDate,
	)

	if err := uc.gateway.Create(ctx, budget); err != nil {
		return nil, err
	}

	return uc.toResponse(budget), nil
}

func (uc *useCase) Update(ctx context.Context, id string, dto *dtos.UpdateBudgetRequest) (*dtos.BudgetResponse, error) {
	budget, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	budget.SetEndDate(dto.EndDate)

	if err := uc.gateway.Update(ctx, budget); err != nil {
		return nil, err
	}

	return uc.toResponse(budget), nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.gateway.Delete(ctx, id)
}

func (uc *useCase) Get(ctx context.Context, id string) (*dtos.BudgetResponse, error) {
	budget, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toResponse(budget), nil
}

func (uc *useCase) List(ctx context.Context, page models.PageRequest) (*models.Page[*dtos.BudgetResponse], error) {
	budgets, count, err := uc.gateway.List(ctx, page)
	if err != nil {
		return nil, err
	}

	responses := make([]*dtos.BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = uc.toResponse(budget)
	}
	return &models.Page[*dtos.BudgetResponse]{
		Page:       page.Page,
		Limit:      page.Limit,
		TotalPages: count / page.Limit,
		Results:    responses,
	}, nil
}

func (uc *useCase) ListByMonth(ctx context.Context, month time.Month, year int) ([]*dtos.BudgetResponse, error) {
	budgets, err := uc.gateway.ListByMonth(ctx, month, year)
	if err != nil {
		return nil, err
	}

	responses := make([]*dtos.BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = uc.toResponse(budget)
	}
	return responses, nil
}

func (uc *useCase) GetSummary(ctx context.Context, month time.Month, year int) (*dtos.BudgetSummaryResponse, error) {
	summary, err := uc.gateway.GetSummary(ctx, month, year)
	if err != nil {
		return nil, err
	}

	return &dtos.BudgetSummaryResponse{
		TotalBudgeted:   summary.TotalBudgeted(),
		TotalSpent:      summary.TotalSpent(),
		Remaining:       summary.Remaining(),
		PercentageSpent: summary.PercentageSpent(),
	}, nil
}

func (uc *useCase) toResponse(budget models.Budget) *dtos.BudgetResponse {
	return &dtos.BudgetResponse{
		ID:          budget.ID(),
		Amount:      budget.Amount(),
		Description: budget.Description(),
		EndDate:     budget.EndDate(),
		Status:      string(budget.Status()),
		CreatedAt:   budget.CreatedAt(),
		UpdatedAt:   budget.UpdatedAt(),
	}
}
