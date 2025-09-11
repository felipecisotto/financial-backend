package factories

import (
	"financial-backend/internal/domain/constants"
	"financial-backend/internal/domain/entities"
	"financial-backend/internal/dto/request"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ExpenseFactory provides factory methods for expense creation
type ExpenseFactory struct{}

// NewExpenseFactory creates a new expense factory
func NewExpenseFactory() *ExpenseFactory {
	return &ExpenseFactory{}
}

// FromCreateRequest creates an expense entity from create request
func (f *ExpenseFactory) FromCreateRequest(req *request.CreateExpenseRequest) *entities.Expense {
	now := time.Now()

	expense := &entities.Expense{
		ID:          uuid.New().String(),
		Description: strings.TrimSpace(req.Description),
		Amount:      normalizeAmount(req.Amount),
		Type:        normalizeExpenseType(req.Type),
		BudgetID:    req.BudgetID,
		Method:      normalizePaymentMethod(req.Method),
		DueDay:      validateDueDay(req.DueDay),
		StartDate:   parseTimeOrNow(req.StartDate),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Handle optional fields
	if req.Recurrency != nil {
		recurrency := normalizeRecurrency(*req.Recurrency)
		expense.Recurrency = &recurrency
	}

	if req.Installments != nil {
		installments := validateInstallments(*req.Installments)
		expense.Installments = &installments
	}

	if req.EndDate != nil {
		endDate := *req.EndDate
		expense.EndDate = &endDate
	}

	return expense
}

// FromUpdateRequest updates an expense entity from update request
func (f *ExpenseFactory) FromUpdateRequest(existing *entities.Expense, req *request.UpdateExpenseRequest) *entities.Expense {
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

	if req.StartDate != nil {
		updated.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		updated.EndDate = req.EndDate
	}

	if req.Installments != nil {
		installments := validateInstallments(*req.Installments)
		updated.Installments = &installments
	}

	return &updated
}

// CreateDefault creates an expense with default values for testing
func (f *ExpenseFactory) CreateDefault() *entities.Expense {
	return f.FromCreateRequest(&request.CreateExpenseRequest{
		Description: "Default Expense",
		Amount:      100.0,
		Type:        constants.ExpenseTypeVariable,
		Method:      constants.PaymentMethodCash,
		DueDay:      15,
		StartDate:   time.Now(),
	})
}

// Helper functions
func normalizeAmount(amount float64) float64 {
	// Round to 2 decimal places
	return float64(int(amount*100+0.5)) / 100
}

func normalizeExpenseType(expenseType string) string {
	normalized := strings.ToUpper(strings.TrimSpace(expenseType))
	// Validate against allowed types
	for _, validType := range constants.ValidExpenseTypes() {
		if normalized == validType {
			return normalized
		}
	}
	return constants.ExpenseTypeVariable // default
}

func normalizePaymentMethod(method string) string {
	normalized := strings.ToUpper(strings.TrimSpace(method))
	for _, validMethod := range constants.ValidPaymentMethods() {
		if normalized == validMethod {
			return normalized
		}
	}
	return constants.PaymentMethodCash // default
}

func normalizeRecurrency(recurrency string) string {
	normalized := strings.ToUpper(strings.TrimSpace(recurrency))
	for _, validRecurrency := range constants.ValidRecurrencyTypes() {
		if normalized == validRecurrency {
			return normalized
		}
	}
	return constants.RecurrencyMonthly // default
}

func validateDueDay(dueDay int) int {
	if dueDay < constants.MinDueDay {
		return constants.MinDueDay
	}
	if dueDay > constants.MaxDueDay {
		return constants.MaxDueDay
	}
	return dueDay
}

func validateInstallments(installments int) int {
	if installments < constants.MinInstallments {
		return constants.MinInstallments
	}
	if installments > constants.MaxInstallments {
		return constants.MaxInstallments
	}
	return installments
}

func parseTimeOrNow(startDate time.Time) time.Time {
	if startDate.IsZero() {
		return time.Now()
	}
	return startDate
}
