package constants

// Expense types
const (
	ExpenseTypeFixed    = "FIXED"
	ExpenseTypeVariable = "VARIABLE"
)

// Payment methods
const (
	PaymentMethodCash         = "CASH"
	PaymentMethodDebitCard    = "DEBIT_CARD"
	PaymentMethodCreditCard   = "CREDIT_CARD"
	PaymentMethodPix          = "PIX"
	PaymentMethodBankTransfer = "BANK_TRANSFER"
)

// Recurrency types
const (
	RecurrencyMonthly = "MONTHLY"
	RecurrencyWeekly  = "WEEKLY"
	RecurrencyYearly  = "YEARLY"
)

// Income sources
const (
	IncomeSourceSalary     = "SALARY"
	IncomeSourceBonus      = "BONUS"
	IncomeSourceInvestment = "INVESTMENT"
	IncomeSourceOther      = "OTHER"
)

// Budget statuses
const (
	BudgetStatusActive    = "ACTIVE"
	BudgetStatusInactive  = "INACTIVE"
	BudgetStatusCompleted = "COMPLETED"
)

// Movement types
const (
	MovementTypeDebit  = "DEBIT"
	MovementTypeCredit = "CREDIT"
	MovementTypeStart  = "START"
)

// Helper functions for validation
func ValidExpenseTypes() []string {
	return []string{ExpenseTypeFixed, ExpenseTypeVariable}
}

func ValidPaymentMethods() []string {
	return []string{
		PaymentMethodCash,
		PaymentMethodDebitCard,
		PaymentMethodCreditCard,
		PaymentMethodPix,
		PaymentMethodBankTransfer,
	}
}

func ValidRecurrencyTypes() []string {
	return []string{RecurrencyMonthly, RecurrencyWeekly, RecurrencyYearly}
}

func ValidIncomeSources() []string {
	return []string{
		IncomeSourceSalary,
		IncomeSourceBonus,
		IncomeSourceInvestment,
		IncomeSourceOther,
	}
}

func ValidBudgetStatuses() []string {
	return []string{BudgetStatusActive, BudgetStatusInactive, BudgetStatusCompleted}
}

func ValidMovementTypes() []string {
	return []string{MovementTypeDebit, MovementTypeCredit, MovementTypeStart}
}