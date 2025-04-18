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

func (c *BudgetMovementController) RegisterRoutes(router *gin.RouterGroup) {
	budgets := router.Group("/movements")
	{
		budgets.POST("", c.Create)
	}
}
