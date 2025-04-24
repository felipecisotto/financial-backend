package controllers

import (
	"financial-backend/internal/dtos"
	budgetmovement "financial-backend/internal/usecases/budget_movement"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BudgetMovementController struct {
	useCase budgetmovement.UseCase
}

func NewBudgetMovementController(uc budgetmovement.UseCase) *BudgetMovementController {
	return &BudgetMovementController{
		useCase: uc,
	}
}

func (c *BudgetMovementController) Create(ctx *gin.Context) {
	var input dtos.BudgetMovementRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.useCase.Create(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *BudgetMovementController) Find(ctx *gin.Context) {
	var params dtos.BudgetMovementParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.useCase.Find(ctx, params)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *BudgetMovementController) ProcessMovements(ctx *gin.Context) {
	if err := c.useCase.CreateRecurrencyMovements(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *BudgetMovementController) RegisterRoutes(router *gin.RouterGroup) {
	budgets := router.Group("/movements")
	{
		budgets.POST("", c.Create)
		budgets.GET("", c.Find)
		budgets.POST("/recurrent", c.ProcessMovements)
	}
}
