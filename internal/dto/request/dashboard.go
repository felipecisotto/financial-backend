package request

// SummaryQueryParams representa os parâmetros para consulta de resumo
type SummaryQueryParams struct {
	Month int `form:"month" binding:"required"`
	Year  int `form:"year" binding:"required"`
}
