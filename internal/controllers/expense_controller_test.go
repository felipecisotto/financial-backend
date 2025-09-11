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

type ExpenseControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockExpenseUseCase
	controller  *ExpenseController
	router      *gin.Engine
}

func (s *ExpenseControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUseCase = &mocks.MockExpenseUseCase{}
	s.controller = NewExpenseController(s.mockUseCase)

	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *ExpenseControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *ExpenseControllerTestSuite) TestCreate_Success() {
	expenseDTO := &dto.ExpenseDTO{
		Description: "Test Expense",
		Amount:      150.0,
		Type:        "VARIABLE",
		Method:      "CREDIT_CARD",
		DueDay:      15,
		StartDate:   time.Now(),
	}

	expectedResponse := &dto.ExpenseResponse{
		ID:           "test-id",
		Description:  expenseDTO.Description,
		Amount:       expenseDTO.Amount,
		Type:         expenseDTO.Type,
		Method:       expenseDTO.Method,
		DueDay:       expenseDTO.DueDay,
		StartDate:    expenseDTO.StartDate,
		EndDate:      expenseDTO.EndDate,
		BudgetID:     expenseDTO.BudgetID,
		Recurrency:   expenseDTO.Recurrency,
		Installments: expenseDTO.Installments,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(dto *dto.ExpenseDTO) bool {
		return dto.Description == "Test Expense" && dto.Amount == 150.0 && dto.Type == "VARIABLE"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(expenseDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response dto.ExpenseResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *ExpenseControllerTestSuite) TestCreate_InvalidJSON() {
	invalidJSON := `{"description":"Test","amount":}`

	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), response["error"], "invalid")
}

func (s *ExpenseControllerTestSuite) TestCreate_MissingRequiredFields() {
	incompleteExpense := map[string]interface{}{
		"description": "Test Expense",
		// Missing required fields: amount, type, method, due_day, start_date
	}

	body, _ := json.Marshal(incompleteExpense)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *ExpenseControllerTestSuite) TestCreate_UseCaseError() {
	expenseDTO := &dto.ExpenseDTO{
		Description: "Test Expense",
		Amount:      150.0,
		Type:        "VARIABLE",
		Method:      "CREDIT_CARD",
		DueDay:      15,
		StartDate:   time.Now(),
	}

	s.mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(dto *dto.ExpenseDTO) bool {
		return dto.Description == "Test Expense" && dto.Amount == 150.0 && dto.Type == "VARIABLE"
	})).Return((*dto.ExpenseResponse)(nil), fmt.Errorf("use case error"))

	body, _ := json.Marshal(expenseDTO)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "use case error", response["error"])
}

func (s *ExpenseControllerTestSuite) TestDelete_Success() {
	expenseID := "test-expense-id"

	s.mockUseCase.On("Delete", mock.Anything, expenseID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
	assert.Empty(s.T(), w.Body.String())
}

func (s *ExpenseControllerTestSuite) TestDelete_UseCaseError() {
	expenseID := "test-expense-id"

	s.mockUseCase.On("Delete", mock.Anything, expenseID).Return(fmt.Errorf("delete error"))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "delete error", response["error"])
}

func (s *ExpenseControllerTestSuite) TestGetByID_Success() {
	expenseID := "test-expense-id"
	expectedResponse := &dto.ExpenseResponse{
		ID:          expenseID,
		Description: "Test Expense",
		Amount:      150.0,
		Type:        "VARIABLE",
		Method:      "CREDIT_CARD",
		DueDay:      15,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.mockUseCase.On("FindByID", mock.Anything, expenseID).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response dto.ExpenseResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.ID, response.ID)
	assert.Equal(s.T(), expectedResponse.Description, response.Description)
}

func (s *ExpenseControllerTestSuite) TestGetByID_NotFound() {
	expenseID := "non-existent-id"

	s.mockUseCase.On("FindByID", mock.Anything, expenseID).Return(nil, fmt.Errorf("expense not found"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/expenses/%s", expenseID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "expense not found", response["error"])
}

func (s *ExpenseControllerTestSuite) TestList_Success() {
	expectedExpenses := []*dto.ExpenseResponse{
		{
			ID:          "expense-1",
			Description: "Expense 1",
			Amount:      100.0,
			Type:        "VARIABLE",
			Method:      "CREDIT_CARD",
			DueDay:      15,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "expense-2",
			Description: "Expense 2",
			Amount:      200.0,
			Type:        "VARIABLE",
			Method:      "DEBIT_CARD",
			DueDay:      20,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	expectedPage := &models.Page[*dto.ExpenseResponse]{
		Page:       1,
		Limit:      10,
		TotalPages: 1,
		Results:    expectedExpenses,
	}

	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(req *dto.ListExpensesRequest) bool {
		return req.Description == "Test" && req.Type == "VARIABLE"
	})).Return(expectedPage, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/expenses?description=Test&type=VARIABLE&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response models.Page[*dto.ExpenseResponse]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), response.Page)
	assert.Equal(s.T(), int64(10), response.Limit)
	assert.Len(s.T(), response.Results, 2)
}

func (s *ExpenseControllerTestSuite) TestList_InvalidQueryParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/expenses?page=invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *ExpenseControllerTestSuite) TestList_UseCaseError() {
	s.mockUseCase.On("List", mock.Anything, mock.MatchedBy(func(req *dto.ListExpensesRequest) bool {
		return true // Accept any request for this error test
	})).Return(nil, fmt.Errorf("list error"))

	req := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "list error", response["error"])
}

func TestExpenseControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseControllerTestSuite))
}
