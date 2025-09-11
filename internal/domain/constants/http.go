package constants

// HTTP status codes (custom if needed)
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
)

// API endpoint patterns
const (
	APIPrefix              = "/api"
	ExpenseEndpoint        = "/expenses"
	IncomeEndpoint         = "/incomes"
	BudgetEndpoint         = "/budgets"
	BudgetMovementEndpoint = "/budget-movements"
	DashboardEndpoint      = "/dashboard"
	HealthEndpoint         = "/ping"
)

// Headers
const (
	HeaderContentType        = "Content-Type"
	HeaderApplicationJSON    = "application/json"
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"
)
