package factories

import (
	"financial-backend/internal/domain/constants"
	"financial-backend/internal/domain/entities"
	"financial-backend/internal/dto/request"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IncomeFactory provides factory methods for income creation
type IncomeFactory struct{}

// NewIncomeFactory creates a new income factory
func NewIncomeFactory() *IncomeFactory {
	return &IncomeFactory{}
}

// FromCreateRequest creates an income entity from create request
func (f *IncomeFactory) FromCreateRequest(req *request.CreateIncomeRequest) *entities.Income {
	now := time.Now()

	income := &entities.Income{
		ID:          uuid.New().String(),
		Description: strings.TrimSpace(req.Description),
		Amount:      normalizeAmount(req.Amount),
		Type:        normalizeIncomeType(req.Type),
		StartDate:   req.StartDate,
		DueDay:      validateDueDay(req.DueDay),
		EndDate:     req.EndDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return income
}

// FromUpdateRequest updates an income entity from update request
func (f *IncomeFactory) FromUpdateRequest(existing *entities.Income, req *request.UpdateIncomeRequest) *entities.Income {
	// Create a copy to avoid modifying the original
	updated := *existing
	updated.UpdatedAt = time.Now()

	// Update only provided fields
	if req.Description != nil {
		updated.Description = strings.TrimSpace(*req.Description)
	}

	if req.Amount != nil {
		updated.Amount = normalizeAmount(*req.Amount)
	}

	if req.Type != nil {
		updated.Type = normalizeIncomeType(*req.Type)
	}

	if req.Date != nil {
		updated.StartDate = *req.Date
	}

	return &updated
}

// CreateDefault creates an income with default values for testing
func (f *IncomeFactory) CreateDefault() *entities.Income {
	return f.FromCreateRequest(&request.CreateIncomeRequest{
		Description: "Default Income",
		Amount:      2000.0,
		Type:        constants.IncomeSourceSalary,
		DueDay:      15,
		StartDate:   time.Now(),
	})
}

// Helper functions
func normalizeIncomeType(incomeType string) string {
	normalized := strings.ToUpper(strings.TrimSpace(incomeType))
	for _, validSource := range constants.ValidIncomeSources() {
		if normalized == validSource {
			return normalized
		}
	}
	return constants.IncomeSourceOther // default
}
