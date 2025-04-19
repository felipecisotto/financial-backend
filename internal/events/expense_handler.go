package events

import (
	"financial-backend/internal/models/events"
	budgetmovement "financial-backend/internal/usecases/budget_movement"
	"financial-backend/pkg/config"

	"gorm.io/gorm"
)

type ExpenseCreatedHandler struct {
	db            *gorm.DB
	createExpense budgetmovement.UseCase
}

func NewExpenseCreatedHandler(db *gorm.DB, createExpense budgetmovement.UseCase) *ExpenseCreatedHandler {
	return &ExpenseCreatedHandler{
		db:            db,
		createExpense: createExpense,
	}
}

func (h *ExpenseCreatedHandler) EventName() string {
	return "ExpenseCreated"
}

func (h *ExpenseCreatedHandler) Handle(e config.Event) {
	event := e.(*events.ExpenseCreatedEvent)
	h.createExpense.CreateExpenseMovement(event.Context, event.Expense)
}
