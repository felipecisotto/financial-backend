package response

// SummaryView representa o resumo financeiro
type SummaryView struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	TotalRemaining float64 `json:"total_remaining"`
}

// SummaryBudgetUtilization representa a utilização do orçamento
type SummaryBudgetUtilization struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Usage       float64 `json:"usage"`
}