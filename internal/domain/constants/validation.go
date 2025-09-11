package constants

// Validation patterns
const (
	UUIDPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	DatePattern = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`
)

// Validation messages
const (
	MsgRequired            = "is required"
	MsgInvalidFormat       = "has invalid format"
	MsgInvalidValue        = "has invalid value"
	MsgTooShort            = "is too short"
	MsgTooLong             = "is too long"
	MsgTooSmall            = "is too small"
	MsgTooLarge            = "is too large"
	MsgInsufficientBudget  = "insufficient budget amount"
	MsgBudgetNotActive     = "budget must be active"
	MsgDateInFuture        = "date cannot be in the future"
	MsgEndDateBeforeStart  = "end date must be after start date"
	MsgInvalidExpenseType  = "invalid expense type"
	MsgInvalidPaymentMethod = "invalid payment method"
	MsgInvalidRecurrency   = "invalid recurrency type"
	MsgInvalidIncomeSource = "invalid income source"
	MsgInvalidBudgetStatus = "invalid budget status"
)