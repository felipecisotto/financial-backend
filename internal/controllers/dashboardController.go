package controllers

import (
	"financial-backend/internal/dtos"
	"financial-backend/internal/usecases/dashboard"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DashboardController struct {
	uc dashboard.UseCase
}

func NewDashboardController(uc dashboard.UseCase) *DashboardController {
	return &DashboardController{
		uc: uc,
	}
}

func (d *DashboardController) GetSummary(ctx *gin.Context) {
	var input dtos.SummaryQueryParams

	if err := ctx.ShouldBindQuery(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := d.uc.GetSummary(ctx, input.Month, input.Year)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (d *DashboardController) SummaryBudgetUsageByMonthYear(ctx *gin.Context) {
	var input dtos.SummaryQueryParams

	if err := ctx.ShouldBindQuery(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := d.uc.SummaryBudgetUsageByMonthYear(ctx, input.Month, input.Year)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

func (d *DashboardController) RegisterRoutes(api *gin.RouterGroup) {
	api = api.Group("/dashboard")
	{
		api.GET("/summary", d.GetSummary)
		api.GET("/budget/utilization", d.SummaryBudgetUsageByMonthYear)
	}
}
