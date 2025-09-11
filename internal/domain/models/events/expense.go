package events

import (
	"context"
	"financial-backend/internal/domain/models"
)

type ExpenseCreatedEvent struct {
	Expense models.Expense
	Context context.Context
}

func (e *ExpenseCreatedEvent) EventName() string {
	return "ExpenseCreated"
}
