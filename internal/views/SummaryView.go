package views

type SummaryView struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	TotalRemaining float64 `json:"total_remaining"`
}
