package factories

import (
	"financial-backend/internal/domain/constants"
	"financial-backend/internal/domain/entities"
	"financial-backend/internal/dto/request"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BudgetMovementFactory provides factory methods for budget movement creation
type BudgetMovementFactory struct{}

// NewBudgetMovementFactory creates a new budget movement factory
func NewBudgetMovementFactory() *BudgetMovementFactory {
	return &BudgetMovementFactory{}
}

// FromCreateRequest creates a budget movement entity from create request
func (f *BudgetMovementFactory) FromCreateRequest(req *request.BudgetMovementRequest) *entities.BudgetMovement {
	now := time.Now()

	movement := &entities.BudgetMovement{
		ID:        uuid.New().String(),
		BudgetId:  req.BudgetId,
		Origin:    strings.TrimSpace(req.Origin),
		Month:     req.Month,
		Year:      req.Year,
		Type:      normalizeMovementType(req.Type),
		Amount:    req.Amount,
		CreatedAt: now,
	}

	return movement
}

// CreateForExpense creates a debit movement for an expense
func (f *BudgetMovementFactory) CreateForExpense(budgetID string, amount int, origin string, month, year int) *entities.BudgetMovement {
	now := time.Now()

	return &entities.BudgetMovement{
		ID:        uuid.New().String(),
		BudgetId:  budgetID,
		Origin:    strings.TrimSpace(origin),
		Month:     month,
		Year:      year,
		Type:      constants.MovementTypeDebit,
		Amount:    amount,
		CreatedAt: now,
	}
}

// CreateForIncome creates a credit movement for an income
func (f *BudgetMovementFactory) CreateForIncome(budgetID string, amount int, origin string, month, year int) *entities.BudgetMovement {
	now := time.Now()

	return &entities.BudgetMovement{
		ID:        uuid.New().String(),
		BudgetId:  budgetID,
		Origin:    strings.TrimSpace(origin),
		Month:     month,
		Year:      year,
		Type:      constants.MovementTypeCredit,
		Amount:    amount,
		CreatedAt: now,
	}
}

// CreateStartMovement creates an initial movement when budget is created
func (f *BudgetMovementFactory) CreateStartMovement(budgetID string, amount int, year int) *entities.BudgetMovement {
	now := time.Now()

	return &entities.BudgetMovement{
		ID:        uuid.New().String(),
		BudgetId:  budgetID,
		Origin:    "Initial budget allocation",
		Month:     int(now.Month()),
		Year:      year,
		Type:      constants.MovementTypeStart,
		Amount:    amount,
		CreatedAt: now,
	}
}

// CreateDefault creates a budget movement with default values for testing
func (f *BudgetMovementFactory) CreateDefault() *entities.BudgetMovement {
	now := time.Now()
	return f.FromCreateRequest(&request.BudgetMovementRequest{
		BudgetId: uuid.New().String(),
		Origin:   "Default Movement",
		Month:    int(now.Month()),
		Year:     now.Year(),
		Type:     constants.MovementTypeDebit,
		Amount:   50,
	})
}

// Helper functions
func normalizeMovementType(movementType string) string {
	normalized := strings.ToUpper(strings.TrimSpace(movementType))
	for _, validType := range constants.ValidMovementTypes() {
		if normalized == validType {
			return normalized
		}
	}
	return constants.MovementTypeDebit // default
}
