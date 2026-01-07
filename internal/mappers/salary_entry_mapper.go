package mappers

import (
	"financial-backend/internal/domain/entities"
	"financial-backend/internal/domain/models"
)

// SalaryEntryModelToEntity converts a salary entry model to an entity
func SalaryEntryModelToEntity(model models.SalaryEntry) *entities.SalaryEntry {
	return &entities.SalaryEntry{
		ID:              model.ID(),
		IncomeID:        model.IncomeID(),
		Month:           model.Month(),
		Year:            model.Year(),
		EntryType:       model.EntryType(),
		Category:        model.Category(),
		CalculationType: model.CalculationType(),
		BaseValue:       model.BaseValue(),
		Multiplier:      model.Multiplier(),
		Amount:          model.Amount(),
		Description:     model.Description(),
		ReferenceDate:   model.ReferenceDate(),
		CreatedAt:       model.CreatedAt(),
		UpdatedAt:       model.UpdatedAt(),
	}
}

// SalaryEntryEntityToModel converts a salary entry entity to a model
func SalaryEntryEntityToModel(entity *entities.SalaryEntry) (models.SalaryEntry, error) {
	return models.NewSalaryEntry(
		entity.ID,
		entity.IncomeID,
		entity.Month,
		entity.Year,
		entity.EntryType,
		entity.Category,
		entity.CalculationType,
		entity.BaseValue,
		entity.Multiplier,
		entity.Amount,
		entity.Description,
		entity.ReferenceDate,
	)
}
