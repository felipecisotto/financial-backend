package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"financial-backend/internal/dto/response"
	"financial-backend/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DashboardControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockDashboardUseCase
	controller  *DashboardController
	router      *gin.Engine
}

func (s *DashboardControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUseCase = &mocks.MockDashboardUseCase{}
	s.controller = NewDashboardController(s.mockUseCase)

	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *DashboardControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

func (s *DashboardControllerTestSuite) TestGetSummary_Success() {
	expectedSummary := response.SummaryView{
		TotalIncome:    5000.0,
		TotalExpense:   3000.0,
		TotalRemaining: 2000.0,
	}

	s.mockUseCase.On("GetSummary", mock.Anything, 1, 2024).Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response response.SummaryView
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedSummary.TotalIncome, response.TotalIncome)
	assert.Equal(s.T(), expectedSummary.TotalExpense, response.TotalExpense)
	assert.Equal(s.T(), expectedSummary.TotalRemaining, response.TotalRemaining)
}

func (s *DashboardControllerTestSuite) TestGetSummary_MissingRequiredParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=1", nil) // missing year
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestGetSummary_InvalidParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=invalid&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestGetSummary_UseCaseError() {
	s.mockUseCase.On("GetSummary", mock.Anything, 1, 2024).Return(response.SummaryView{}, fmt.Errorf("summary error"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestSummaryBudgetUsageByMonthYear_Success() {
	expectedUtilizations := []response.SummaryBudgetUtilization{
		{
			Description: "Food Budget",
			Amount:      500.0,
			Usage:       300.0,
		},
		{
			Description: "Transport Budget",
			Amount:      200.0,
			Usage:       150.0,
		},
	}

	s.mockUseCase.On("SummaryBudgetUsageByMonthYear", mock.Anything, 1, 2024).Return(expectedUtilizations, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/budget/utilization?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response []response.SummaryBudgetUtilization
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), response, 2)
	assert.Equal(s.T(), "Food Budget", response[0].Description)
	assert.Equal(s.T(), 500.0, response[0].Amount)
	assert.Equal(s.T(), 300.0, response[0].Usage)
}

func (s *DashboardControllerTestSuite) TestSummaryBudgetUsageByMonthYear_MissingRequiredParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/budget/utilization?month=1", nil) // missing year
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestSummaryBudgetUsageByMonthYear_InvalidParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/budget/utilization?month=invalid&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestSummaryBudgetUsageByMonthYear_UseCaseError() {
	s.mockUseCase.On("SummaryBudgetUsageByMonthYear", mock.Anything, 1, 2024).Return(nil, fmt.Errorf("utilization error"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/budget/utilization?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestGetMonthlyEvolution_Success() {
	expectedEvolution := response.MonthlyEvolutionView{
		{Month: 1, Income: 5000.0, Expense: 3000.0},
		{Month: 2, Income: 5200.0, Expense: 3200.0},
		{Month: 3, Income: 5000.0, Expense: 2800.0},
		{Month: 4, Income: 5100.0, Expense: 3100.0},
		{Month: 5, Income: 5300.0, Expense: 3300.0},
		{Month: 6, Income: 5200.0, Expense: 3000.0},
		{Month: 7, Income: 5400.0, Expense: 3400.0},
		{Month: 8, Income: 5500.0, Expense: 3500.0},
		{Month: 9, Income: 5300.0, Expense: 3200.0},
		{Month: 10, Income: 5600.0, Expense: 3600.0},
		{Month: 11, Income: 5700.0, Expense: 3700.0},
		{Month: 12, Income: 5800.0, Expense: 3800.0},
	}

	s.mockUseCase.On("GetMonthlyEvolution", mock.Anything, 2024).Return(expectedEvolution, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly-evolution?year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response response.MonthlyEvolutionView
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), response, 12)
	assert.Equal(s.T(), 1, response[0].Month)
	assert.Equal(s.T(), 5000.0, response[0].Income)
	assert.Equal(s.T(), 3000.0, response[0].Expense)
	assert.Equal(s.T(), 12, response[11].Month)
	assert.Equal(s.T(), 5800.0, response[11].Income)
	assert.Equal(s.T(), 3800.0, response[11].Expense)
}

func (s *DashboardControllerTestSuite) TestGetMonthlyEvolution_MissingRequiredParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly-evolution", nil) // missing year
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestGetMonthlyEvolution_InvalidParams() {
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly-evolution?year=invalid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestGetMonthlyEvolution_UseCaseError() {
	s.mockUseCase.On("GetMonthlyEvolution", mock.Anything, 2024).Return(response.MonthlyEvolutionView(nil), fmt.Errorf("evolution error"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly-evolution?year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func TestDashboardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardControllerTestSuite))
}
