package mappers

import (
	"financial-backend/internal/entities"
	"financial-backend/internal/models"
)

func ToExpenseModel(entity *entities.Expense) models.Expense {
	var budget *models.Budget
	if entity.Budget != nil {
		newBudget := ToBudgetModel(entity.Budget)
		budget = &newBudget
	}
	expense, _ := models.NewExpense(
		entity.ID,
		entity.Description,
		entity.Amount,
		entity.Type,
		entity.BudgetID,
		entity.Recurrency,
		entity.Method,
		entity.Installments,
		entity.DueDay,
		entity.StartDate,
		entity.EndDate,
		budget,
	)
	return expense
}
