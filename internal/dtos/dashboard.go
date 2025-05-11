package dtos

type SummaryQueryParams struct {
	Month int `form:"month" binding:"required"`
	Year  int `form:"year" binding:"required"`
}
