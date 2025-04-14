package controllers

import (
	"net/http"

	"financial-backend/internal/dtos"
	"financial-backend/internal/usecases/expense"

	"github.com/gin-gonic/gin"
)

type ExpenseController struct {
	UseCase expense.UseCase
}

func NewExpenseController(useCase expense.UseCase) *ExpenseController {
	return &ExpenseController{
		UseCase: useCase,
	}
}

func (c *ExpenseController) Create(ctx *gin.Context) {
	var input dtos.ExpenseDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.UseCase.Create(ctx, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ExpenseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.UseCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ExpenseController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	response, err := c.UseCase.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ExpenseController) List(ctx *gin.Context) {
	var input dtos.ListExpensesRequest
	if err := ctx.ShouldBindQuery(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.UseCase.List(ctx, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ExpenseController) RegisterRoutes(router *gin.RouterGroup) {
	expenses := router.Group("/expenses")
	{
		expenses.POST("", c.Create)
		expenses.DELETE("/:id", c.Delete)
		expenses.GET("/:id", c.GetByID)
		expenses.GET("", c.List)
	}
}
