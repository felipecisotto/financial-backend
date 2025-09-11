package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto"
	"financial-backend/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IncomeControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockIncomeUseCase
	controller  *IncomeController
	router      *gin.Engine
}

func (s *IncomeControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUseCase = &mocks.MockIncomeUseCase{}
	s.controller = NewIncomeController(s.mockUseCase)

	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *IncomeControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *IncomeControllerTestSuite) TestCreate_Success() {
	createReq := &dto.CreateIncomeRequest{
		Description: "Test Income",
		Amount:      1000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
	}

	expectedResponse := &dto.IncomeResponse{
		ID:          "test-id",
		Description: "Test Income",
		Amount:      1000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(dto *dto.CreateIncomeRequest) bool {
		return dto.Description == "Test Income" && dto.Amount == 1000.0 && dto.Type == "FIXED"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/incomes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response dto.IncomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *IncomeControllerTestSuite) TestCreate_InvalidJSON() {
	invalidJSON := `{"description":"Test","amount":}`

	req := httptest.NewRequest(http.MethodPost, "/api/incomes", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *IncomeControllerTestSuite) TestCreate_MissingRequiredFields() {
	incompleteIncome := map[string]interface{}{
		"description": "Test Income",
		// Missing required fields: amount, type, due_day, start_date
	}

	body, _ := json.Marshal(incompleteIncome)
	req := httptest.NewRequest(http.MethodPost, "/api/incomes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *IncomeControllerTestSuite) TestCreate_UseCaseError() {
	createReq := &dto.CreateIncomeRequest{
		Description: "Test Income",
		Amount:      1000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(dto *dto.CreateIncomeRequest) bool {
		return dto.Description == "Test Income" && dto.Amount == 1000.0 && dto.Type == "FIXED"
	})).Return(nil, fmt.Errorf("create error"))

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/incomes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *IncomeControllerTestSuite) TestUpdate_Success() {
	incomeID := "test-income-id"
	amount := 1500.0
	updateReq := &dto.UpdateIncomeRequest{
		Amount: &amount,
	}

	expectedResponse := &dto.IncomeResponse{
		ID:     incomeID,
		Amount: 1500.0,
	}

	s.mockUseCase.On("Update", mock.Anything, incomeID, mock.MatchedBy(func(dto *dto.UpdateIncomeRequest) bool {
		return dto.Amount != nil && *dto.Amount == 1500.0
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/incomes/%s", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response dto.IncomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Amount, response.Amount)
}

func (s *IncomeControllerTestSuite) TestUpdate_InvalidJSON() {
	incomeID := "test-income-id"
	invalidJSON := `{"amount":}`

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/incomes/%s", incomeID), bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *IncomeControllerTestSuite) TestUpdate_UseCaseError() {
	incomeID := "test-income-id"
	amount := 1500.0
	updateReq := &dto.UpdateIncomeRequest{
		Amount: &amount,
	}

	s.mockUseCase.On("Update", mock.Anything, incomeID, mock.MatchedBy(func(req *dto.UpdateIncomeRequest) bool {
		return req.Amount != nil && *req.Amount == 1500.0
	})).Return(nil, fmt.Errorf("update error"))

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/incomes/%s", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *IncomeControllerTestSuite) TestDelete_Success() {
	incomeID := "test-income-id"

	s.mockUseCase.On("Delete", mock.Anything, incomeID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/incomes/%s", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	assert.Empty(s.T(), w.Body.String())
}

func (s *IncomeControllerTestSuite) TestDelete_UseCaseError() {
	incomeID := "test-income-id"

	s.mockUseCase.On("Delete", mock.Anything, incomeID).Return(fmt.Errorf("delete error"))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/incomes/%s", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *IncomeControllerTestSuite) TestGet_Success() {
	incomeID := "test-income-id"
	expectedResponse := &dto.IncomeResponse{
		ID:          incomeID,
		Description: "Test Income",
		Amount:      1000.0,
		Type:        "FIXED",
		DueDay:      1,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockUseCase.On("Get", mock.Anything, incomeID).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/incomes/%s", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response dto.IncomeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *IncomeControllerTestSuite) TestGet_NotFound() {
	incomeID := "non-existent-id"

	s.mockUseCase.On("Get", mock.Anything, incomeID).Return(nil, fmt.Errorf("income not found"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/incomes/%s", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *IncomeControllerTestSuite) TestList_Success() {
	expectedIncomes := []*dto.IncomeResponse{
		{
			ID:          "income-1",
			Description: "Income 1",
			Amount:      1000.0,
			Type:        "FIXED",
			DueDay:      1,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "income-2",
			Description: "Income 2",
			Amount:      1500.0,
			Type:        "VARIABLE",
			DueDay:      5,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	expectedPage := &models.Page[*dto.IncomeResponse]{
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Results:    expectedIncomes,
	}

	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(params dto.ListIncomeParams) bool {
		return params.Type == "FIXED" && params.Description == "Test"
	})).Return(expectedPage, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/incomes?type=FIXED&description=Test&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response models.Page[*dto.IncomeResponse]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), response.Page)
	assert.Equal(s.T(), int64(10), response.Limit)
	assert.Len(s.T(), response.Results, 2)
}

func (s *IncomeControllerTestSuite) TestList_InvalidQueryParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/incomes?page=invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *IncomeControllerTestSuite) TestList_UseCaseError() {
	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(params dto.ListIncomeParams) bool {
		return true // Accept any params for this error test
	})).Return(nil, fmt.Errorf("list error"))

	req := httptest.NewRequest(http.MethodGet, "/api/incomes", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func TestIncomeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeControllerTestSuite))
}
