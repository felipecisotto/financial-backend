package response

import "time"

// SalaryEntryResponse representa a estrutura de resposta para uma entrada de salário
type SalaryEntryResponse struct {
	ID              string     `json:"id"`
	IncomeID        string     `json:"income_id"`
	Month           int        `json:"month"`
	Year            int        `json:"year"`
	EntryType       string     `json:"entry_type"`
	Category        string     `json:"category"`
	CalculationType string     `json:"calculation_type"`
	BaseValue       float64    `json:"base_value"`
	Multiplier      float64    `json:"multiplier"`
	Amount          float64    `json:"amount"`
	Description     *string    `json:"description,omitempty"`
	ReferenceDate   *time.Time `json:"reference_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// SalarySummaryResponse representa a estrutura de resposta para o resumo de salário
type SalarySummaryResponse struct {
	IncomeID       string                `json:"income_id"`
	Description    string                `json:"description"`
	Month          int                   `json:"month"`
	Year           int                   `json:"year"`
	BaseAmount     float64               `json:"base_amount"`
	Entries        []SalaryEntryResponse `json:"entries"`
	TotalAdditions float64               `json:"total_additions"`
	TotalDiscounts float64               `json:"total_discounts"`
	NetAmount      float64               `json:"net_amount"`
}

// MonthlySummary representa o resumo mensal de salário
type MonthlySummary struct {
	Month          int     `json:"month"`
	BaseAmount     float64 `json:"base_amount"`
	TotalAdditions float64 `json:"total_additions"`
	TotalDiscounts float64 `json:"total_discounts"`
	NetAmount      float64 `json:"net_amount"`
}

// SalaryHistoryResponse representa a estrutura de resposta para o histórico de salário
type SalaryHistoryResponse struct {
	IncomeID         string           `json:"income_id"`
	Description      string           `json:"description"`
	Year             int              `json:"year"`
	MonthlySummaries []MonthlySummary `json:"monthly_summaries"`
}

// DashboardSalarySummaryResponse representa a estrutura de resposta para o resumo de salário do dashboard
type DashboardSalarySummaryResponse struct {
	Month    int                     `json:"month"`
	Year     int                     `json:"year"`
	Salaries []SalarySummaryResponse `json:"salaries"`
}
