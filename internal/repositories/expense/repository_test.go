package expense

import (
	"context"
	"testing"
	"time"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"
	"financial-backend/internal/repositories/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExpenseRepositoryTestSuite struct {
	suite.Suite
	repository Repository
	cleanup    func()
}

func (s *ExpenseRepositoryTestSuite) SetupSuite() {
	db, cleanup := testutils.SetupTestDatabase(s.T())
	s.cleanup = cleanup
	s.repository = NewRepository(db)
}

func (s *ExpenseRepositoryTestSuite) TearDownSuite() {
	s.cleanup()
}

func (s *ExpenseRepositoryTestSuite) TestCreate_Success() {
	expense := &entities.Expense{
		ID:          uuid.New().String(),
		Description: "Test Expense",
		Amount:      100.50,
		Type:        "VARIABLE",
		Method:      "PIX",
		DueDay:      15,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repository.Create(context.Background(), expense)
	assert.NoError(s.T(), err)

	// Verify the expense was created
	savedExpense, err := s.repository.Get(context.Background(), expense.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expense.Description, savedExpense.Description)
	assert.Equal(s.T(), expense.Amount, savedExpense.Amount)
	assert.Equal(s.T(), expense.Type, savedExpense.Type)
	assert.Equal(s.T(), expense.Method, savedExpense.Method)
}

func (s *ExpenseRepositoryTestSuite) TestCreate_WithBudgetID() {
	budgetID := uuid.New().String()
	expense := &entities.Expense{
		ID:          uuid.New().String(),
		Description: "Expense with Budget",
		Amount:      200.00,
		Type:        "VARIABLE",
		BudgetID:    &budgetID,
		Method:      "CREDIT_CARD",
		DueDay:      20,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repository.Create(context.Background(), expense)
	assert.NoError(s.T(), err)

	// Verify the expense was created with budget ID
	savedExpense, err := s.repository.Get(context.Background(), expense.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), savedExpense.BudgetID)
	assert.Equal(s.T(), budgetID, *savedExpense.BudgetID)
}

func (s *ExpenseRepositoryTestSuite) TestGet_NotFound() {
	nonExistentID := uuid.New().String()
	
	expense, err := s.repository.Get(context.Background(), nonExistentID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), expense)
}

func (s *ExpenseRepositoryTestSuite) TestUpdate_Success() {
	// First create an expense
	expense := &entities.Expense{
		ID:          uuid.New().String(),
		Description: "Original Description",
		Amount:      100.00,
		Type:        "VARIABLE",
		Method:      "PIX",
		DueDay:      15,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repository.Create(context.Background(), expense)
	assert.NoError(s.T(), err)

	// Update the expense
	expense.Description = "Updated Description"
	expense.Amount = 150.00
	expense.UpdatedAt = time.Now()

	err = s.repository.Update(context.Background(), expense)
	assert.NoError(s.T(), err)

	// Verify the update
	updatedExpense, err := s.repository.Get(context.Background(), expense.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Description", updatedExpense.Description)
	assert.Equal(s.T(), 150.00, updatedExpense.Amount)
}

func (s *ExpenseRepositoryTestSuite) TestDelete_Success() {
	// First create an expense
	expense := &entities.Expense{
		ID:          uuid.New().String(),
		Description: "To be deleted",
		Amount:      100.00,
		Type:        "VARIABLE",
		Method:      "PIX",
		DueDay:      15,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repository.Create(context.Background(), expense)
	assert.NoError(s.T(), err)

	// Delete the expense
	err = s.repository.Delete(context.Background(), expense.ID)
	assert.NoError(s.T(), err)

	// Verify the expense is deleted
	deletedExpense, err := s.repository.Get(context.Background(), expense.ID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), deletedExpense)
}

func (s *ExpenseRepositoryTestSuite) TestList_WithFilters() {
	ctx := context.Background()
	
	// Create test expenses
	expenses := []*entities.Expense{
		{
			ID:          uuid.New().String(),
			Description: "Food expense",
			Amount:      50.00,
			Type:        "VARIABLE",
			Method:      "PIX",
			DueDay:      15,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Description: "Transport expense",
			Amount:      30.00,
			Type:        "FIXED",
			Method:      "CREDIT_CARD",
			DueDay:      20,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Description: "Another food expense",
			Amount:      75.00,
			Type:        "VARIABLE",
			Method:      "PIX",
			DueDay:      15,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, expense := range expenses {
		err := s.repository.Create(ctx, expense)
		assert.NoError(s.T(), err)
	}

	// Test list with description filter
	page := models.PageRequest{Page: 1, Limit: 10}
	results, count, err := s.repository.List(ctx, "food", "", "", "", "", "", page)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count) // Should find 2 food expenses
	assert.Len(s.T(), results, 2)

	// Test list with type filter
	results, count, err = s.repository.List(ctx, "", "FIXED", "", "", "", "", page)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count) // Should find 1 fixed expense
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "FIXED", results[0].Type)

	// Test list with method filter
	results, count, err = s.repository.List(ctx, "", "", "", "", "", "PIX", page)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count) // Should find 2 PIX expenses
	assert.Len(s.T(), results, 2)
}

func (s *ExpenseRepositoryTestSuite) TestList_WithPagination() {
	ctx := context.Background()
	
	// Create multiple test expenses
	for i := 0; i < 5; i++ {
		expense := &entities.Expense{
			ID:          uuid.New().String(),
			Description: "Test expense " + string(rune(i+'1')),
			Amount:      float64(100 + i*10),
			Type:        "VARIABLE",
			Method:      "PIX",
			DueDay:      15,
			StartDate:   time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := s.repository.Create(ctx, expense)
		assert.NoError(s.T(), err)
	}

	// Test pagination - first page
	page1 := models.PageRequest{Page: 1, Limit: 2}
	results, count, err := s.repository.List(ctx, "", "", "", "", "", "", page1)
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), count, int64(4)) // Should have at least 5 expenses
	assert.LessOrEqual(s.T(), len(results), 2) // Should return max 2 results

	// Test pagination - second page
	page2 := models.PageRequest{Page: 2, Limit: 2}
	results2, count2, err := s.repository.List(ctx, "", "", "", "", "", "", page2)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), count, count2) // Same total count
	assert.LessOrEqual(s.T(), len(results2), 2)

	// Verify results from different pages are different (if we have enough data)
	if len(results) > 0 && len(results2) > 0 {
		assert.NotEqual(s.T(), results[0].ID, results2[0].ID)
	}
}

func (s *ExpenseRepositoryTestSuite) TestSummaryByMonth() {
	ctx := context.Background()
	
	// Create expenses for a specific month
	year, month := 2024, 1
	testDate := time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC)
	
	expenses := []*entities.Expense{
		{
			ID:          uuid.New().String(),
			Description: "January expense 1",
			Amount:      100.00,
			Type:        "VARIABLE",
			Method:      "PIX",
			DueDay:      15,
			StartDate:   testDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Description: "January expense 2",
			Amount:      200.00,
			Type:        "VARIABLE",
			Method:      "PIX",
			DueDay:      20,
			StartDate:   testDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, expense := range expenses {
		err := s.repository.Create(ctx, expense)
		assert.NoError(s.T(), err)
	}

	// Test summary for the month
	summary, err := s.repository.SummaryByMonth(ctx, month, year)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 300.00, summary) // 100 + 200
}

func TestExpenseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseRepositoryTestSuite))
}