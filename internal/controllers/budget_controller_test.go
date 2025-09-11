package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"financial-backend/internal/dto"
	"financial-backend/internal/domain/models"
	"financial-backend/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BudgetControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockBudgetUseCase
	controller  *BudgetController
	router      *gin.Engine
}

func (s *BudgetControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	
	s.mockUseCase = &mocks.MockBudgetUseCase{}
	s.controller = NewBudgetController(s.mockUseCase)
	
	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *BudgetControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *BudgetControllerTestSuite) TestCreate_Success() {
	createReq := dto.CreateBudgetRequest{
		Description: "Test Budget",
		Amount:      1000.0,
	}

	expectedResponse := dto.BudgetResponse{
		ID:          "test-id",
		Description: "Test Budget",
		Amount:      1000.0,
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(dto dto.CreateBudgetRequest) bool {
		return dto.Description == "Test Budget" && dto.Amount == 1000.0
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	var response dto.BudgetResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *BudgetControllerTestSuite) TestCreate_InvalidJSON() {
	invalidJSON := `{"description":"Test","amount":}`
	
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetControllerTestSuite) TestCreate_MissingRequiredFields() {
	incompleteBudget := map[string]interface{}{
		"description": "Test Budget",
		// Missing required field: amount
	}

	body, _ := json.Marshal(incompleteBudget)
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetControllerTestSuite) TestCreate_UseCaseError() {
	createReq := dto.CreateBudgetRequest{
		Description: "Test Budget",
		Amount:      1000.0,
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(req dto.CreateBudgetRequest) bool {
		return req.Description == "Test Budget" && req.Amount == 1000.0
	})).Return(dto.BudgetResponse{}, fmt.Errorf("create error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *BudgetControllerTestSuite) TestUpdate_Success() {
	budgetID := "test-budget-id"
	endDate := time.Now().AddDate(0, 1, 0) // 1 month from now
	updateReq := dto.UpdateBudgetRequest{
		EndDate: endDate,
	}

	expectedResponse := dto.BudgetResponse{
		ID:      budgetID,
		EndDate: &endDate,
	}

	s.mockUseCase.On("Update", mock.Anything, budgetID, mock.MatchedBy(func(req *dto.UpdateBudgetRequest) bool {
		return req.EndDate.Equal(endDate)
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/budgets/%s", budgetID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *BudgetControllerTestSuite) TestUpdate_InvalidJSON() {
	budgetID := "test-budget-id"
	invalidJSON := `{"end_date":}`
	
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/budgets/%s", budgetID), bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetControllerTestSuite) TestDelete_Success() {
	budgetID := "test-budget-id"
	
	s.mockUseCase.On("Delete", mock.Anything, budgetID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/budgets/%s", budgetID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	assert.Empty(s.T(), w.Body.String())
}

func (s *BudgetControllerTestSuite) TestDelete_UseCaseError() {
	budgetID := "test-budget-id"
	
	s.mockUseCase.On("Delete", mock.Anything, budgetID).Return(fmt.Errorf("delete error"))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/budgets/%s", budgetID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *BudgetControllerTestSuite) TestGet_Success() {
	budgetID := "test-budget-id"
	expectedResponse := dto.BudgetResponse{
		ID:          budgetID,
		Description: "Test Budget",
		Amount:      1000.0,
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockUseCase.On("Get", mock.Anything, budgetID).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/budgets/%s", budgetID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response dto.BudgetResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *BudgetControllerTestSuite) TestGet_NotFound() {
	budgetID := "non-existent-id"
	
	s.mockUseCase.On("Get", mock.Anything, budgetID).Return(dto.BudgetResponse{}, fmt.Errorf("budget not found"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/budgets/%s", budgetID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *BudgetControllerTestSuite) TestList_Success() {
	expectedBudgets := []dto.BudgetResponse{
		{
			ID:          "budget-1",
			Description: "Budget 1",
			Amount:      1000.0,
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "budget-2",
			Description: "Budget 2",
			Amount:      1500.0,
			Status:      "ACTIVE",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	expectedPage := &models.Page[dto.BudgetResponse]{
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Results:    expectedBudgets,
	}

	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(params dto.BudgetListParams) bool {
		return params.Description == "Test" && params.Status == "ACTIVE"
	})).Return(expectedPage, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/budgets?description=Test&status=ACTIVE&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response models.Page[dto.BudgetResponse]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), response.Page)
	assert.Equal(s.T(), int64(10), response.Limit)
	assert.Len(s.T(), response.Results, 2)
}

func (s *BudgetControllerTestSuite) TestList_InvalidQueryParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/budgets?page=invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *BudgetControllerTestSuite) TestList_UseCaseError() {
	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(params dto.BudgetListParams) bool {
		return true // Accept any params for this error test
	})).Return(nil, fmt.Errorf("list error"))

	req := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func TestBudgetControllerTestSuite(t *testing.T) {
	suite.Run(t, new(BudgetControllerTestSuite))
}