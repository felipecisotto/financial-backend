package constants

import "time"

// Business limits and rules
const (
	// Amount limits
	MinAmount        = 0.01
	MaxExpenseAmount = 1000000.0
	MaxIncomeAmount  = 10000000.0
	MaxBudgetAmount  = 1000000.0

	// Date constraints
	MinBudgetDuration = 24 * time.Hour
	MaxFutureDays     = 1

	// Installments
	MinInstallments = 1
	MaxInstallments = 120

	// Description constraints
	MinDescriptionLength = 3
	MaxDescriptionLength = 100

	// Due day constraints
	MinDueDay = 1
	MaxDueDay = 31

	// Pagination
	DefaultPageSize = 20
	MaxPageSize     = 100
	MinPageSize     = 1
)