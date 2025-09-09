package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"financial-backend/internal/views"
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
	expectedSummary := views.SummaryView{
		TotalIncome:    5000.0,
		TotalExpense:   3000.0,
		TotalRemaining: 2000.0,
	}

	s.mockUseCase.On("GetSummary", mock.Anything, 1, 2024).Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	var response views.SummaryView
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
	s.mockUseCase.On("GetSummary", mock.Anything, 1, 2024).Return(views.SummaryView{}, fmt.Errorf("summary error"))

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary?month=1&year=2024", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *DashboardControllerTestSuite) TestSummaryBudgetUsageByMonthYear_Success() {
	expectedUtilizations := []views.SummaryBudgetUtilization{
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
	
	var response []views.SummaryBudgetUtilization
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

func TestDashboardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardControllerTestSuite))
}