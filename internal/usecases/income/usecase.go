package income

import (
	"context"
	"financial-backend/internal/dto"
	"financial-backend/internal/gateways"
	"financial-backend/internal/domain/models"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, dto *dto.CreateIncomeRequest) (*dto.IncomeResponse, error)
	Update(ctx context.Context, id string, dto *dto.UpdateIncomeRequest) (*dto.IncomeResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*dto.IncomeResponse, error)
	List(ctx context.Context, params dto.ListIncomeParams) (*models.Page[*dto.IncomeResponse], error)
}

type useCase struct {
	gateway gateways.IncomeGateway
}

func NewUseCase(gateway gateways.IncomeGateway) UseCase {
	return &useCase{gateway: gateway}
}

func (uc *useCase) Create(ctx context.Context, dto *dto.CreateIncomeRequest) (*dto.IncomeResponse, error) {
	income, err := models.NewIncome(
		uuid.New().String(),
		dto.Description,
		dto.Amount,
		models.IncomeType(dto.Type),
		dto.DueDay,
		dto.StartDate,
		dto.EndDate,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.gateway.Create(ctx, income); err != nil {
		return nil, err
	}

	return uc.toResponse(income), nil
}

func (uc *useCase) Update(ctx context.Context, id string, dto *dto.UpdateIncomeRequest) (*dto.IncomeResponse, error) {
	income, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// if dto.Description != nil {
	// 	income.SetDescription(*dto.Description)
	// }
	// if dto.Amount != nil {
	// 	income.SetAmount(*dto.Amount)
	// }
	// if dto.Date != nil {
	// 	income.SetDate(dto.Date)
	// }
	// income.SetUpdatedAt(time.Now())

	if err := uc.gateway.Update(ctx, income); err != nil {
		return nil, err
	}

	return uc.toResponse(income), nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.gateway.Delete(ctx, id)
}

func (uc *useCase) Get(ctx context.Context, id string) (*dto.IncomeResponse, error) {
	income, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toResponse(income), nil
}

func (uc *useCase) List(ctx context.Context, params dto.ListIncomeParams) (*models.Page[*dto.IncomeResponse], error) {
	fmt.Printf("params %v", params)
	incomes, count, err := uc.gateway.List(ctx, params.Type, params.Description, models.PageRequest{
		Limit: params.Limit,
		Page:  params.Page,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.IncomeResponse, len(incomes))
	for i, income := range incomes {
		responses[i] = uc.toResponse(income)
	}

	return &models.Page[*dto.IncomeResponse]{
		Results:    responses,
		TotalPages: int64(math.Ceil(float64(count) / float64(params.Limit))),
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (uc *useCase) toResponse(income models.Income) *dto.IncomeResponse {
	if income == nil {
		return nil
	}

	return &dto.IncomeResponse{
		ID:          income.ID(),
		Description: income.Description(),
		Amount:      income.Amount(),
		Type:        string(income.Type()),
		DueDay:      income.DueDay(),
		StartDate:   income.StartDate(),
		EndDate:     income.EndDate(),
		CreatedAt:   income.CreatedAt(),
		UpdatedAt:   income.UpdatedAt(),
	}
}
