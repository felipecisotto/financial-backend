package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"financial-backend/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SalaryControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.MockSalaryUseCase
	controller  *SalaryController
	router      *gin.Engine
}

func (s *SalaryControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockUseCase = &mocks.MockSalaryUseCase{}
	s.controller = NewSalaryController(s.mockUseCase)

	s.router = gin.New()
	apiGroup := s.router.Group("/api")
	s.controller.RegisterRoutes(apiGroup)
}

func (s *SalaryControllerTestSuite) TearDownTest() {
	s.mockUseCase.AssertExpectations(s.T())
}

// ============================================================================
// POST /api/salary/:income_id/activate Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestActivateTracking_Success_CLTMode() {
	incomeID := "test-income-id"
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "CLT",
		HourlyRate:   50.0,
	}

	s.mockUseCase.On("ActivateTracking", mock.Anything, incomeID, mock.MatchedBy(func(req request.ActivateSalaryTrackingRequest) bool {
		return req.DiscountMode == "CLT" && req.HourlyRate == 50.0
	})).Return(nil)

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Salary tracking activated successfully", response["message"])
}

func (s *SalaryControllerTestSuite) TestActivateTracking_Success_CustomMode() {
	incomeID := "test-income-id"
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "custom",
		HourlyRate:   75.5,
	}

	s.mockUseCase.On("ActivateTracking", mock.Anything, incomeID, mock.MatchedBy(func(req request.ActivateSalaryTrackingRequest) bool {
		return req.DiscountMode == "custom" && req.HourlyRate == 75.5
	})).Return(nil)

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *SalaryControllerTestSuite) TestActivateTracking_InvalidJSON() {
	incomeID := "test-income-id"
	invalidJSON := `{"discount_mode":"CLT","hourly_rate":}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestActivateTracking_InvalidDiscountMode() {
	incomeID := "test-income-id"
	activateReq := map[string]interface{}{
		"discount_mode": "INVALID",
		"hourly_rate":   50.0,
	}

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestActivateTracking_MissingHourlyRate() {
	incomeID := "test-income-id"
	activateReq := map[string]interface{}{
		"discount_mode": "CLT",
	}

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestActivateTracking_NegativeHourlyRate() {
	incomeID := "test-income-id"
	activateReq := map[string]interface{}{
		"discount_mode": "CLT",
		"hourly_rate":   -50.0,
	}

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestActivateTracking_IncomeNotFound() {
	incomeID := "non-existent-id"
	activateReq := request.ActivateSalaryTrackingRequest{
		DiscountMode: "CLT",
		HourlyRate:   50.0,
	}

	s.mockUseCase.On("ActivateTracking", mock.Anything, incomeID, mock.Anything).Return(fmt.Errorf("income not found"))

	body, _ := json.Marshal(activateReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/activate", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// POST /api/salary/:income_id/extra-hours Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestAddExtraHours_Success_WithAllFields() {
	incomeID := "test-income-id"
	rate := 100.0
	referenceDate := "2024-01-15"
	extraHoursReq := request.AddExtraHoursRequest{
		Month:         1,
		Year:          2024,
		Hours:         10.5,
		Rate:          &rate,
		ReferenceDate: &referenceDate,
	}

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Month:          1,
		Year:           2024,
		BaseAmount:     5000.0,
		TotalAdditions: 1050.0,
		TotalDiscounts: 500.0,
		NetAmount:      5550.0,
	}

	s.mockUseCase.On("AddExtraHours", mock.Anything, incomeID, mock.MatchedBy(func(req request.AddExtraHoursRequest) bool {
		return req.Month == 1 && req.Year == 2024 && req.Hours == 10.5 && req.Rate != nil && *req.Rate == 100.0
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.IncomeID, response.IncomeID)
	assert.Equal(s.T(), expectedResponse.NetAmount, response.NetAmount)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_Success_WithoutRate() {
	incomeID := "test-income-id"
	extraHoursReq := request.AddExtraHoursRequest{
		Month: 1,
		Year:  2024,
		Hours: 8.0,
	}

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Month:          1,
		Year:           2024,
		BaseAmount:     5000.0,
		TotalAdditions: 800.0,
		TotalDiscounts: 500.0,
		NetAmount:      5300.0,
	}

	s.mockUseCase.On("AddExtraHours", mock.Anything, incomeID, mock.MatchedBy(func(req request.AddExtraHoursRequest) bool {
		return req.Month == 1 && req.Year == 2024 && req.Hours == 8.0 && req.Rate == nil
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.IncomeID, response.IncomeID)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_InvalidHours_Zero() {
	incomeID := "test-income-id"
	extraHoursReq := map[string]interface{}{
		"month": 1,
		"year":  2024,
		"hours": 0,
	}

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_InvalidHours_Negative() {
	incomeID := "test-income-id"
	extraHoursReq := map[string]interface{}{
		"month": 1,
		"year":  2024,
		"hours": -5.0,
	}

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_InvalidMonth() {
	incomeID := "test-income-id"
	extraHoursReq := map[string]interface{}{
		"month": 13,
		"year":  2024,
		"hours": 8.0,
	}

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_InvalidYear() {
	incomeID := "test-income-id"
	extraHoursReq := map[string]interface{}{
		"month": 1,
		"year":  2019,
		"hours": 8.0,
	}

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddExtraHours_UseCaseError() {
	incomeID := "test-income-id"
	extraHoursReq := request.AddExtraHoursRequest{
		Month: 1,
		Year:  2024,
		Hours: 8.0,
	}

	s.mockUseCase.On("AddExtraHours", mock.Anything, incomeID, mock.Anything).Return(nil, fmt.Errorf("failed to add extra hours"))

	body, _ := json.Marshal(extraHoursReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/extra-hours", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// POST /api/salary/:income_id/on-call Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestAddOnCall_Success() {
	incomeID := "test-income-id"
	referenceDate := "2024-01-20"
	onCallReq := request.AddOnCallRequest{
		Month:         1,
		Year:          2024,
		Hours:         12.0,
		ReferenceDate: &referenceDate,
	}

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Month:          1,
		Year:           2024,
		BaseAmount:     5000.0,
		TotalAdditions: 1200.0,
		TotalDiscounts: 500.0,
		NetAmount:      5700.0,
	}

	s.mockUseCase.On("AddOnCall", mock.Anything, incomeID, mock.MatchedBy(func(req request.AddOnCallRequest) bool {
		return req.Month == 1 && req.Year == 2024 && req.Hours == 12.0
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(onCallReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.IncomeID, response.IncomeID)
	assert.Equal(s.T(), expectedResponse.NetAmount, response.NetAmount)
}

func (s *SalaryControllerTestSuite) TestAddOnCall_InvalidJSON() {
	incomeID := "test-income-id"
	invalidJSON := `{"month":1,"year":2024,"hours":}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddOnCall_InvalidHours_Zero() {
	incomeID := "test-income-id"
	onCallReq := map[string]interface{}{
		"month": 1,
		"year":  2024,
		"hours": 0,
	}

	body, _ := json.Marshal(onCallReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddOnCall_InvalidHours_Negative() {
	incomeID := "test-income-id"
	onCallReq := map[string]interface{}{
		"month": 1,
		"year":  2024,
		"hours": -10.0,
	}

	body, _ := json.Marshal(onCallReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddOnCall_InvalidMonth() {
	incomeID := "test-income-id"
	onCallReq := map[string]interface{}{
		"month": 0,
		"year":  2024,
		"hours": 12.0,
	}

	body, _ := json.Marshal(onCallReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddOnCall_UseCaseError() {
	incomeID := "test-income-id"
	onCallReq := request.AddOnCallRequest{
		Month: 1,
		Year:  2024,
		Hours: 12.0,
	}

	s.mockUseCase.On("AddOnCall", mock.Anything, incomeID, mock.Anything).Return(nil, fmt.Errorf("failed to add on-call"))

	body, _ := json.Marshal(onCallReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/on-call", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// GET /api/salary/:income_id/summary/:month/:year Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestGetSummary_Success() {
	incomeID := "test-income-id"
	month := 1
	year := 2024

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Description:    "Test Salary",
		Month:          month,
		Year:           year,
		BaseAmount:     5000.0,
		TotalAdditions: 1000.0,
		TotalDiscounts: 800.0,
		NetAmount:      5200.0,
	}

	s.mockUseCase.On("GetSummary", mock.Anything, incomeID, month, year).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/summary/%d/%d", incomeID, month, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.IncomeID, response.IncomeID)
	assert.Equal(s.T(), expectedResponse.NetAmount, response.NetAmount)
}

func (s *SalaryControllerTestSuite) TestGetSummary_NoEntries_ReturnsBaseAmount() {
	incomeID := "test-income-id"
	month := 2
	year := 2024

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Description:    "Test Salary",
		Month:          month,
		Year:           year,
		BaseAmount:     5000.0,
		TotalAdditions: 0.0,
		TotalDiscounts: 0.0,
		NetAmount:      5000.0,
		Entries:        []response.SalaryEntryResponse{},
	}

	s.mockUseCase.On("GetSummary", mock.Anything, incomeID, month, year).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/summary/%d/%d", incomeID, month, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 5000.0, response.NetAmount)
	assert.Equal(s.T(), 0.0, response.TotalAdditions)
	assert.Equal(s.T(), 0.0, response.TotalDiscounts)
}

func (s *SalaryControllerTestSuite) TestGetSummary_InvalidMonth() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/summary/invalid/2024", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Invalid month parameter", response["error"])
}

func (s *SalaryControllerTestSuite) TestGetSummary_InvalidYear() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/summary/1/invalid", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Invalid year parameter", response["error"])
}

func (s *SalaryControllerTestSuite) TestGetSummary_UseCaseError() {
	incomeID := "test-income-id"
	month := 1
	year := 2024

	s.mockUseCase.On("GetSummary", mock.Anything, incomeID, month, year).Return(nil, fmt.Errorf("failed to get summary"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/summary/%d/%d", incomeID, month, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// DELETE /api/salary/entries/:entry_id Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestRemoveEntry_Success() {
	entryID := "test-entry-id"

	s.mockUseCase.On("RemoveEntry", mock.Anything, entryID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/salary/entries/%s", entryID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Entry removed successfully", response["message"])
}

func (s *SalaryControllerTestSuite) TestRemoveEntry_NotFound() {
	entryID := "non-existent-entry-id"

	s.mockUseCase.On("RemoveEntry", mock.Anything, entryID).Return(fmt.Errorf("entry not found"))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/salary/entries/%s", entryID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func (s *SalaryControllerTestSuite) TestRemoveEntry_UseCaseError() {
	entryID := "test-entry-id"

	s.mockUseCase.On("RemoveEntry", mock.Anything, entryID).Return(fmt.Errorf("database error"))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/salary/entries/%s", entryID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// GET /api/salary/:income_id/history?year=2024 Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestGetHistory_Success_Returns12Months() {
	incomeID := "test-income-id"
	year := 2024

	monthlySummaries := make([]response.MonthlySummary, 12)
	for i := 0; i < 12; i++ {
		monthlySummaries[i] = response.MonthlySummary{
			Month:          i + 1,
			BaseAmount:     5000.0,
			TotalAdditions: 500.0,
			TotalDiscounts: 300.0,
			NetAmount:      5200.0,
		}
	}

	expectedResponse := &response.SalaryHistoryResponse{
		IncomeID:         incomeID,
		Description:      "Test Salary",
		Year:             year,
		MonthlySummaries: monthlySummaries,
	}

	s.mockUseCase.On("GetHistory", mock.Anything, incomeID, year).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history?year=%d", incomeID, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response response.SalaryHistoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.IncomeID, response.IncomeID)
	assert.Equal(s.T(), expectedResponse.Year, response.Year)
	assert.Len(s.T(), response.MonthlySummaries, 12)
}

func (s *SalaryControllerTestSuite) TestGetHistory_Success_CorrectValues() {
	incomeID := "test-income-id"
	year := 2024

	monthlySummaries := []response.MonthlySummary{
		{Month: 1, BaseAmount: 5000.0, TotalAdditions: 1000.0, TotalDiscounts: 500.0, NetAmount: 5500.0},
		{Month: 2, BaseAmount: 5000.0, TotalAdditions: 800.0, TotalDiscounts: 400.0, NetAmount: 5400.0},
		{Month: 3, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 4, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 5, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 6, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 7, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 8, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 9, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 10, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 11, BaseAmount: 5000.0, TotalAdditions: 0.0, TotalDiscounts: 0.0, NetAmount: 5000.0},
		{Month: 12, BaseAmount: 5000.0, TotalAdditions: 1500.0, TotalDiscounts: 700.0, NetAmount: 5800.0},
	}

	expectedResponse := &response.SalaryHistoryResponse{
		IncomeID:         incomeID,
		Description:      "Test Salary",
		Year:             year,
		MonthlySummaries: monthlySummaries,
	}

	s.mockUseCase.On("GetHistory", mock.Anything, incomeID, year).Return(expectedResponse, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history?year=%d", incomeID, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response response.SalaryHistoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 5500.0, response.MonthlySummaries[0].NetAmount)
	assert.Equal(s.T(), 5400.0, response.MonthlySummaries[1].NetAmount)
	assert.Equal(s.T(), 5800.0, response.MonthlySummaries[11].NetAmount)
}

func (s *SalaryControllerTestSuite) TestGetHistory_InvalidYear_Missing() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestGetHistory_InvalidYear_TooOld() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history?year=2019", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestGetHistory_InvalidYear_NotNumeric() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history?year=invalid", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestGetHistory_UseCaseError() {
	incomeID := "test-income-id"
	year := 2024

	s.mockUseCase.On("GetHistory", mock.Anything, incomeID, year).Return(nil, fmt.Errorf("failed to get history"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/salary/%s/history?year=%d", incomeID, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// ============================================================================
// POST /api/salary/:income_id/entries Tests (Generic Entry)
// ============================================================================

func (s *SalaryControllerTestSuite) TestAddEntry_Success() {
	incomeID := "test-income-id"
	description := "Bonus"
	entryReq := request.CreateSalaryEntryRequest{
		Month:           1,
		Year:            2024,
		EntryType:       "addition",
		Category:        "bonus",
		CalculationType: "fixed",
		BaseValue:       1000.0,
		Multiplier:      1.0,
		Description:     &description,
	}

	expectedResponse := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Month:          1,
		Year:           2024,
		BaseAmount:     5000.0,
		TotalAdditions: 1000.0,
		TotalDiscounts: 500.0,
		NetAmount:      5500.0,
	}

	s.mockUseCase.On("AddEntry", mock.Anything, incomeID, mock.MatchedBy(func(req request.CreateSalaryEntryRequest) bool {
		return req.Month == 1 && req.Year == 2024 && req.EntryType == "addition" && req.BaseValue == 1000.0
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(entryReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/entries", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	var response response.SalarySummaryResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResponse.NetAmount, response.NetAmount)
}

func (s *SalaryControllerTestSuite) TestAddEntry_InvalidEntryType() {
	incomeID := "test-income-id"
	entryReq := map[string]interface{}{
		"month":            1,
		"year":             2024,
		"entry_type":       "invalid",
		"category":         "bonus",
		"calculation_type": "fixed",
		"base_value":       1000.0,
		"multiplier":       1.0,
	}

	body, _ := json.Marshal(entryReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/entries", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestAddEntry_InvalidCalculationType() {
	incomeID := "test-income-id"
	entryReq := map[string]interface{}{
		"month":            1,
		"year":             2024,
		"entry_type":       "addition",
		"category":         "bonus",
		"calculation_type": "invalid",
		"base_value":       1000.0,
		"multiplier":       1.0,
	}

	body, _ := json.Marshal(entryReq)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/entries", incomeID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ============================================================================
// POST /api/salary/:income_id/recalculate/:month/:year Tests
// ============================================================================

func (s *SalaryControllerTestSuite) TestRecalculateCLT_Success() {
	incomeID := "test-income-id"
	month := 1
	year := 2024

	s.mockUseCase.On("RecalculateCLT", mock.Anything, incomeID, month, year).Return(nil)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/recalculate/%d/%d", incomeID, month, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "CLT entries recalculated successfully", response["message"])
}

func (s *SalaryControllerTestSuite) TestRecalculateCLT_InvalidMonth() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/recalculate/invalid/2024", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestRecalculateCLT_InvalidYear() {
	incomeID := "test-income-id"

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/recalculate/1/invalid", incomeID), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *SalaryControllerTestSuite) TestRecalculateCLT_UseCaseError() {
	incomeID := "test-income-id"
	month := 1
	year := 2024

	s.mockUseCase.On("RecalculateCLT", mock.Anything, incomeID, month, year).Return(fmt.Errorf("recalculation failed"))

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/salary/%s/recalculate/%d/%d", incomeID, month, year), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

func TestSalaryControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SalaryControllerTestSuite))
}
