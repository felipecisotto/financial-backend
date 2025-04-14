package controllers

import (
	"net/http"

	"financial-backend/internal/dtos"
	"financial-backend/internal/usecases/budget"

	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	useCase budget.UseCase
}

func NewBudgetController(useCase budget.UseCase) *BudgetController {
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
	var params dtos.BudgetListParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros inválidos"})
		return
	}

	response, err := c.useCase.List(ctx, params)
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
	}
}
