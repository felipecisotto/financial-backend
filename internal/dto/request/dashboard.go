package request

// SummaryQueryParams representa os parâmetros para consulta de resumo
type SummaryQueryParams struct {
	Month int `form:"month" binding:"required"`
	Year  int `form:"year" binding:"required"`
}

// MonthlyEvolutionQueryParams representa os parâmetros para consulta de evolução mensal
type MonthlyEvolutionQueryParams struct {
	Year int `form:"year" binding:"required"`
}
