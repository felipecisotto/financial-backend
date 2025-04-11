package controllers

import (
	"net/http"

	"financial-backend/internal/dtos"
	"financial-backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type IncomeController struct {
	UseCase usecases.IncomeUseCase
}

func NewIncomeController(useCase usecases.IncomeUseCase) *IncomeController {
	return &IncomeController{UseCase: useCase}
}

func (c *IncomeController) Create(ctx *gin.Context) {
	var req dtos.CreateIncomeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.UseCase.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *IncomeController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.UpdateIncomeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.UseCase.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *IncomeController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.UseCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *IncomeController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	response, err := c.UseCase.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *IncomeController) List(ctx *gin.Context) {
	responses, err := c.UseCase.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert []*IncomeResponse to []IncomeResponse
	incomes := make([]dtos.IncomeResponse, len(responses))
	for i, r := range responses {
		incomes[i] = *r
	}

	ctx.JSON(http.StatusOK, dtos.ListIncomesResponse{
		Incomes: incomes,
		Total:   int64(len(responses)),
	})
}

func (c *IncomeController) ListByType(ctx *gin.Context) {
	incomeType := ctx.Param("type")
	responses, err := c.UseCase.ListByType(ctx, incomeType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert []*IncomeResponse to []IncomeResponse
	incomes := make([]dtos.IncomeResponse, len(responses))
	for i, r := range responses {
		incomes[i] = *r
	}

	ctx.JSON(http.StatusOK, dtos.ListIncomesResponse{
		Incomes: incomes,
		Total:   int64(len(responses)),
	})
}

func (c *IncomeController) RegisterRoutes(router *gin.RouterGroup) {
	incomes := router.Group("/incomes")
	{
		incomes.POST("", c.Create)
		incomes.PUT("/:id", c.Update)
		incomes.DELETE("/:id", c.Delete)
		incomes.GET("/:id", c.Get)
		incomes.GET("", c.List)
		incomes.GET("/type/:type", c.ListByType)
	}
}
