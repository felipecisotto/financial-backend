package views

type SummaryBudgetUtilization struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Usage       float64 `json:"usage"`
}
