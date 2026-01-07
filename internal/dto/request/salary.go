package request

// ActivateSalaryTrackingRequest representa a requisição para ativar o rastreamento de salário
type ActivateSalaryTrackingRequest struct {
	DiscountMode string  `json:"discount_mode" binding:"required,oneof=CLT custom"`
	HourlyRate   float64 `json:"hourly_rate" binding:"required,gt=0"`
}

// CreateSalaryEntryRequest representa a requisição para criar uma entrada de salário
type CreateSalaryEntryRequest struct {
	Month           int     `json:"month" binding:"required,min=1,max=12"`
	Year            int     `json:"year" binding:"required,min=2020"`
	EntryType       string  `json:"entry_type" binding:"required,oneof=discount addition"`
	Category        string  `json:"category" binding:"required"`
	CalculationType string  `json:"calculation_type" binding:"required,oneof=fixed percentage computed"`
	BaseValue       float64 `json:"base_value" binding:"required"`
	Multiplier      float64 `json:"multiplier" binding:"required"`
	Description     *string `json:"description,omitempty"`
	ReferenceDate   *string `json:"reference_date,omitempty"` // ISO format YYYY-MM-DD
}

// AddExtraHoursRequest representa a requisição para adicionar horas extras
type AddExtraHoursRequest struct {
	Month         int      `json:"month" binding:"required,min=1,max=12"`
	Year          int      `json:"year" binding:"required,min=2020"`
	Hours         float64  `json:"hours" binding:"required,gt=0"`
	Rate          *float64 `json:"rate,omitempty"`
	ReferenceDate *string  `json:"reference_date,omitempty"`
}

// AddOnCallRequest representa a requisição para adicionar plantão
type AddOnCallRequest struct {
	Month         int     `json:"month" binding:"required,min=1,max=12"`
	Year          int     `json:"year" binding:"required,min=2020"`
	Hours         float64 `json:"hours" binding:"required,gt=0"`
	ReferenceDate *string `json:"reference_date,omitempty"`
}

// SalaryHistoryQueryParams representa os parâmetros de query para histórico de salário
type SalaryHistoryQueryParams struct {
	Year int `form:"year" binding:"required,min=2020"`
}

// SalarySummaryQueryParams representa os parâmetros de query para resumo de salário
type SalarySummaryQueryParams struct {
	Month int `form:"month" binding:"required,min=1,max=12"`
	Year  int `form:"year" binding:"required,min=2020"`
}
