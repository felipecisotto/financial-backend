package controllers

import (
	"net/http"
	"strconv"

	"financial-backend/internal/dto/request"
	"financial-backend/internal/usecases/salary"

	"github.com/gin-gonic/gin"
)

type SalaryController struct {
	uc salary.UseCase
}

func NewSalaryController(uc salary.UseCase) *SalaryController {
	return &SalaryController{
		uc: uc,
	}
}

// ActivateTracking activates salary tracking for an income
// POST /api/salary/:income_id/activate
func (c *SalaryController) ActivateTracking(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	var req request.ActivateSalaryTrackingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.uc.ActivateTracking(ctx, incomeID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Salary tracking activated successfully"})
}

// AddEntry adds a generic salary entry
// POST /api/salary/:income_id/entries
func (c *SalaryController) AddEntry(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	var req request.CreateSalaryEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := c.uc.AddEntry(ctx, incomeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, summary)
}

// AddExtraHours adds extra hours to the salary
// POST /api/salary/:income_id/extra-hours
func (c *SalaryController) AddExtraHours(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	var req request.AddExtraHoursRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := c.uc.AddExtraHours(ctx, incomeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, summary)
}

// AddOnCall adds on-call hours to the salary
// POST /api/salary/:income_id/on-call
func (c *SalaryController) AddOnCall(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	var req request.AddOnCallRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := c.uc.AddOnCall(ctx, incomeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, summary)
}

// RemoveEntry removes a salary entry by ID
// DELETE /api/salary/entries/:entry_id
func (c *SalaryController) RemoveEntry(ctx *gin.Context) {
	entryID := ctx.Param("entry_id")

	if err := c.uc.RemoveEntry(ctx, entryID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry removed successfully"})
}

// GetSummary retrieves the monthly salary summary
// GET /api/salary/:income_id/summary/:month/:year
func (c *SalaryController) GetSummary(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	month, err := strconv.Atoi(ctx.Param("month"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month parameter"})
		return
	}

	year, err := strconv.Atoi(ctx.Param("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}

	summary, err := c.uc.GetSummary(ctx, incomeID, month, year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetHistory retrieves the annual salary history
// GET /api/salary/:income_id/history
func (c *SalaryController) GetHistory(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	var params request.SalaryHistoryQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history, err := c.uc.GetHistory(ctx, incomeID, params.Year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, history)
}

// RecalculateCLT recalculates CLT entries for a specific month
// POST /api/salary/:income_id/recalculate/:month/:year
func (c *SalaryController) RecalculateCLT(ctx *gin.Context) {
	incomeID := ctx.Param("income_id")

	month, err := strconv.Atoi(ctx.Param("month"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month parameter"})
		return
	}

	year, err := strconv.Atoi(ctx.Param("year"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
		return
	}

	if err := c.uc.RecalculateCLT(ctx, incomeID, month, year); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "CLT entries recalculated successfully"})
}

// RegisterRoutes registers all salary routes with the provided router group
func (c *SalaryController) RegisterRoutes(api *gin.RouterGroup) {
	salary := api.Group("/salary")
	{
		salary.POST("/:income_id/activate", c.ActivateTracking)
		salary.POST("/:income_id/entries", c.AddEntry)
		salary.POST("/:income_id/extra-hours", c.AddExtraHours)
		salary.POST("/:income_id/on-call", c.AddOnCall)
		salary.DELETE("/entries/:entry_id", c.RemoveEntry)
		salary.GET("/:income_id/summary/:month/:year", c.GetSummary)
		salary.GET("/:income_id/history", c.GetHistory)
		salary.POST("/:income_id/recalculate/:month/:year", c.RecalculateCLT)
	}
}
