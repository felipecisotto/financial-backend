package income

import (
	"context"
	"financial-backend/internal/dtos"
	"financial-backend/internal/gateways"
	"financial-backend/internal/models"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, dto *dtos.CreateIncomeRequest) (*dtos.IncomeResponse, error)
	Update(ctx context.Context, id string, dto *dtos.UpdateIncomeRequest) (*dtos.IncomeResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*dtos.IncomeResponse, error)
	List(ctx context.Context) ([]*dtos.IncomeResponse, error)
	ListByType(ctx context.Context, incomeType string) ([]*dtos.IncomeResponse, error)
}

type useCase struct {
	gateway gateways.IncomeGateway
}

func NewUseCase(gateway gateways.IncomeGateway) UseCase {
	return &useCase{gateway: gateway}
}

func (uc *useCase) Create(ctx context.Context, dto *dtos.CreateIncomeRequest) (*dtos.IncomeResponse, error) {
	income := models.NewIncome(
		uuid.New().String(),
		dto.Description,
		dto.Amount,
		models.IncomeType(dto.Type),
		&dto.DueDay,
	)

	if err := uc.gateway.Create(ctx, income); err != nil {
		return nil, err
	}

	return uc.toResponse(income), nil
}

func (uc *useCase) Update(ctx context.Context, id string, dto *dtos.UpdateIncomeRequest) (*dtos.IncomeResponse, error) {
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

func (uc *useCase) Get(ctx context.Context, id string) (*dtos.IncomeResponse, error) {
	income, err := uc.gateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toResponse(income), nil
}

func (uc *useCase) List(ctx context.Context) ([]*dtos.IncomeResponse, error) {
	incomes, err := uc.gateway.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dtos.IncomeResponse, len(incomes))
	for i, income := range incomes {
		responses[i] = uc.toResponse(income)
	}
	return responses, nil
}

func (uc *useCase) ListByType(ctx context.Context, incomeType string) ([]*dtos.IncomeResponse, error) {
	incomes, err := uc.gateway.ListByType(ctx, models.IncomeType(incomeType))
	if err != nil {
		return nil, err
	}

	responses := make([]*dtos.IncomeResponse, len(incomes))
	for i, income := range incomes {
		responses[i] = uc.toResponse(income)
	}
	return responses, nil
}


func (uc *useCase) toResponse(income models.Income) *dtos.IncomeResponse {
	response := &dtos.IncomeResponse{
		ID:          income.ID(),
		Description: income.Description(),
		Amount:      income.Amount(),
		Type:        string(income.Type()),
		CreatedAt:   income.CreatedAt(),
		UpdatedAt:   income.UpdatedAt(),
	}

	return response
}
