package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"financial-backend/internal/dtos"
	"financial-backend/internal/models"
	"financial-backend/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BudgetMovementControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockBudgetMovementUseCase
	controller  *BudgetMovementController
	router      *gin.Engine
}

func (s *BudgetMovementControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	
	s.mockUseCase = &mocks.MockBudgetMovementUseCase{}
	s.controller = NewBudgetMovementController(s.mockUseCase)
	
	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *BudgetMovementControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *BudgetMovementControllerTestSuite) TestCreate_Success() {
	createReq := dtos.BudgetMovementRequest{
		BudgetId: "budget-id",
		Origin:   "MANUAL",
		Month:    1,
		Year:     2024,
		Type:     "DEBIT",
		Amount:   100,
	}

	expectedResponse := dtos.BudgetMovementResponse{
		ID:     "movement-id",
		Origin: "MANUAL",
		Month:  1,
		Year:   2024,
		Type:   "DEBIT",
		Amount: 100,
		Budget: dtos.BudgetResponse{
			ID:          "budget-id",
			Description: "Test Budget",
		},
		CreatedAt: time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(req dtos.BudgetMovementRequest) bool {
		return req.BudgetId == "budget-id" && req.Amount == 100
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/movements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	var response dtos.BudgetMovementResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Amount, response.Amount)
}

func (s *BudgetMovementControllerTestSuite) TestCreate_InvalidJSON() {
	invalidJSON := `{"budget_id":"test","amount":}`
	
	req := httptest.NewRequest(http.MethodPost, "/api/movements", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetMovementControllerTestSuite) TestCreate_UseCaseError() {
	createReq := dtos.BudgetMovementRequest{
		BudgetId: "budget-id",
		Origin:   "MANUAL",
		Amount:   100,
	}

	s.mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("dtos.BudgetMovementRequest")).Return(dtos.BudgetMovementResponse{}, fmt.Errorf("create error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/movements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *BudgetMovementControllerTestSuite) TestFind_Success() {
	expectedMovements := []dtos.BudgetMovementResponse{
		{
			ID:     "movement-1",
			Origin: "EXPENSE",
			Month:  1,
			Year:   2024,
			Amount: 150,
		},
		{
			ID:     "movement-2",
			Origin: "MANUAL",
			Month:  1,
			Year:   2024,
			Amount: 200,
		},
	}

	expectedPage := models.Page[dtos.BudgetMovementResponse]{
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Results:    expectedMovements,
	}

	s.mockUseCase.On("Find", mock.Anything, mock.MatchedBy(func(params dtos.BudgetMovementParams) bool {
		return params.Month == 1 && params.Year == 2024
	})).Return(expectedPage, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/movements?month=1&year=2024&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response models.Page[dtos.BudgetMovementResponse]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), response.Page)
	assert.Equal(s.T(), int64(10), response.Limit)
	assert.Len(s.T(), response.Results, 2)
}

func (s *BudgetMovementControllerTestSuite) TestFind_InvalidQueryParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/movements?month=invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetMovementControllerTestSuite) TestFind_UseCaseError() {
	s.mockUseCase.On("Find", mock.Anything, mock.AnythingOfType("dtos.BudgetMovementParams")).Return(models.Page[dtos.BudgetMovementResponse]{}, fmt.Errorf("find error"))

	req := httptest.NewRequest(http.MethodGet, "/api/movements", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *BudgetMovementControllerTestSuite) TestProcessMovements_Success() {
	s.mockUseCase.On("CreateRecurrencyMovements", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/movements/recurrent", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	assert.Empty(s.T(), w.Body.String())
}

func (s *BudgetMovementControllerTestSuite) TestProcessMovements_UseCaseError() {
	s.mockUseCase.On("CreateRecurrencyMovements", mock.Anything).Return(fmt.Errorf("process error"))

	req := httptest.NewRequest(http.MethodPost, "/api/movements/recurrent", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func TestBudgetMovementControllerTestSuite(t *testing.T) {
	suite.Run(t, new(BudgetMovementControllerTestSuite))
}