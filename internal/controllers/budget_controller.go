package controllers

import (
	"net/http"
	"strconv"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"financial-backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	useCase usecases.BudgetUseCase
}

func NewBudgetController(useCase usecases.BudgetUseCase) *BudgetController {
	return &BudgetController{useCase: useCase}
}

func (c *BudgetController) Create(ctx *gin.Context) {
	var input dtos.CreateBudgetRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.useCase.Create(ctx, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *BudgetController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dtos.UpdateBudgetRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.useCase.Update(ctx, id, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.useCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *BudgetController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	response, err := c.useCase.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetController) List(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if(err != nil) {
		limit = 50
	}
	page, err := strconv.Atoi(ctx.Query("page"))

	if(err != nil) {
		page = 1
	}
	
	pageRequest := models.PageRequest{
		Limit: int64(limit),
		Page:  int64(page),
	}
	response, err := c.useCase.List(ctx, pageRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetController) ListByMonth(ctx *gin.Context) {
	monthStr := ctx.Query("month")
	yearStr := ctx.Query("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	response, err := c.useCase.ListByMonth(ctx, time.Month(month), year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetController) GetSummary(ctx *gin.Context) {
	monthStr := ctx.Query("month")
	yearStr := ctx.Query("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	response, err := c.useCase.GetSummary(ctx, time.Month(month), year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetController) RegisterRoutes(router *gin.RouterGroup) {
	budgets := router.Group("/budgets")
	{
		budgets.POST("", c.Create)
		budgets.PUT("/:id", c.Update)
		budgets.DELETE("/:id", c.Delete)
		budgets.GET("/:id", c.Get)
		budgets.GET("", c.List)
		budgets.GET("/month", c.ListByMonth)
		budgets.GET("/summary", c.GetSummary)
	}
}
