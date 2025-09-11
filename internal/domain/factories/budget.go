package factories

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"financial-backend/internal/domain/entities"
	"financial-backend/internal/dto/request"
)

// BudgetFactory provides factory methods for budget creation
type BudgetFactory struct{}

// NewBudgetFactory creates a new budget factory
func NewBudgetFactory() *BudgetFactory {
	return &BudgetFactory{}
}

// FromCreateRequest creates a budget entity from create request
func (f *BudgetFactory) FromCreateRequest(req *request.CreateBudgetRequest) *entities.Budget {
	now := time.Now()

	budget := &entities.Budget{
		ID:          uuid.New().String(),
		Description: strings.TrimSpace(req.Description),
		Amount:      normalizeAmount(req.Amount),
		EndDate:     req.EndDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return budget
}

// FromUpdateRequest updates a budget entity from update request
func (f *BudgetFactory) FromUpdateRequest(existing *entities.Budget, req *request.UpdateBudgetRequest) *entities.Budget {
	// Create a copy to avoid modifying the original
	updated := *existing
	updated.UpdatedAt = time.Now()

	// Update the end date
	updated.EndDate = &req.EndDate

	return &updated
}

// CreateDefault creates a budget with default values for testing
func (f *BudgetFactory) CreateDefault() *entities.Budget {
	now := time.Now()
	endDate := now.AddDate(0, 1, 0) // 1 month from now

	return f.FromCreateRequest(&request.CreateBudgetRequest{
		Description: "Default Budget",
		Amount:      1000.0,
		EndDate:     &endDate,
	})
}
