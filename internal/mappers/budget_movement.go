package mappers

import (
	"financial-backend/internal/dtos"
	"financial-backend/internal/entities"
	"financial-backend/internal/models"

	"github.com/google/uuid"
)

// ToEntity converts a BudgetMovement model to a BudgetMovement entity
func ToBudgetMovementEntity(bm models.BudgetMovement) entities.BudgetMovement {
	return entities.BudgetMovement{
		ID:                bm.ID(),
		BudgetId:          bm.BudgetId(),
		Origin:            bm.Origin(),
		Month:             bm.Month(),
		Year:              bm.Year(),
		Type:              string(bm.Type()),
		Amount:            bm.Amount(),
		CreatedAt:         bm.CreatedAt(),
		OriginDescription: nil,
	}
}

// ToModel converts a BudgetMovement entity to a BudgetMovement model
func ToBudgetMovementModel(bmEntity entities.BudgetMovement) models.BudgetMovement {
	budget := ToBudgetModel(&bmEntity.Budget)
	return models.NewBudgetMovement(
		bmEntity.ID,
		bmEntity.BudgetId,
		budget,
		bmEntity.Origin,
		bmEntity.OriginDescription,
		bmEntity.Month,
		bmEntity.Year,
		models.MovementType(bmEntity.Type),
		bmEntity.Amount,
	)
}

// ToDTO converts a BudgetMovement model to a BudgetMovementResponse DTO
func ToBudgetMovementDTO(bm models.BudgetMovement) dtos.BudgetMovementResponse {
	if bm.Budget() != nil {
		return dtos.BudgetMovementResponse{
			ID:                bm.ID(),
			Origin:            bm.Origin(),
			OriginDescription: bm.OriginDescription(),
			Budget:            ToBudgetResponse(bm.Budget()),
			Month:             bm.Month(),
			Year:              bm.Year(),
			Type:              string(bm.Type()),
			Amount:            bm.Amount(),
			CreatedAt:         bm.CreatedAt(),
		}
	} else {
		return dtos.BudgetMovementResponse{
			ID:                bm.ID(),
			Origin:            bm.Origin(),
			OriginDescription: bm.OriginDescription(),
			Month:             bm.Month(),
			Year:              bm.Year(),
			Type:              string(bm.Type()),
			Amount:            bm.Amount(),
			CreatedAt:         bm.CreatedAt(),
		}
	}

}

func FromDTOToBudgetMovementModel(bm dtos.BudgetMovementRequest) models.BudgetMovement {
	return models.NewBudgetMovement(
		uuid.New().String(),
		bm.BudgetId,
		nil,
		bm.Origin,
		nil,
		bm.Month,
		bm.Year,
		models.MovementType(bm.Type),
		bm.Amount,
	)
}
