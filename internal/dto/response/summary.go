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

// MonthlyData representa os dados financeiros de um mês
type MonthlyData struct {
	Month   int     `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// MonthlyEvolutionView representa a evolução mensal de receitas e despesas
type MonthlyEvolutionView []MonthlyData
