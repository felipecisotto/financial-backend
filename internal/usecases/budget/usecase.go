package budget

import (
	"context"
	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"financial-backend/internal/gateways"
	"financial-backend/internal/mappers"
	"math"
)

type UseCase interface {
	Create(ctx context.Context, dto request.CreateBudgetRequest) (response.BudgetResponse, error)
	Update(ctx context.Context, id string, dto *request.UpdateBudgetRequest) (response.BudgetResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (response.BudgetResponse, error)
	List(ctx context.Context, params request.BudgetListParams) (*models.Page[response.BudgetResponse], error)
}

type useCase struct {
	gateway gateways.BudgetGateway
}

func NewUseCase(gateway gateways.BudgetGateway) UseCase {
	return &useCase{gateway: gateway}
}

func (uc *useCase) Create(ctx context.Context, dto request.CreateBudgetRequest) (response.BudgetResponse, error) {
	budget := mappers.FromDTOToBudgetModel(dto)

	if err := uc.gateway.Create(ctx, budget); err != nil {
		return response.BudgetResponse{}, err
	}

	return mappers.ToBudgetResponse(budget), nil
}

func (uc *useCase) Update(ctx context.Context, id string, dto *request.UpdateBudgetRequest) (response.BudgetResponse, error) {
	budget, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return response.BudgetResponse{}, err
	}

	budget.SetEndDate(dto.EndDate)

	if err := uc.gateway.Update(ctx, budget); err != nil {
		return response.BudgetResponse{}, err
	}

	return mappers.ToBudgetResponse(budget), nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.gateway.Delete(ctx, id)
}

func (uc *useCase) Get(ctx context.Context, id string) (response.BudgetResponse, error) {
	budget, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return response.BudgetResponse{}, err
	}
	return mappers.ToBudgetResponse(budget), nil
}

func (uc *useCase) List(ctx context.Context, dto request.BudgetListParams) (*models.Page[response.BudgetResponse], error) {
	budgets, count, err := uc.gateway.List(ctx, dto.Status, dto.Description, models.PageRequest{
		Page:  dto.Page,
		Limit: dto.Limit,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]response.BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = mappers.ToBudgetResponse(budget)
	}
	return &models.Page[response.BudgetResponse]{
		Page:       dto.Page,
		Limit:      dto.Limit,
		TotalPages: int64(math.Ceil(float64(count) / float64(dto.Limit))),
		Results:    responses,
	}, nil
}
