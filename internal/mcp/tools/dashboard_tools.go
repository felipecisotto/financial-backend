package tools

import (
	"context"
	"fmt"
	"time"

	dashboardUC "financial-backend/internal/usecases/dashboard"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DashboardTools provides MCP tools for dashboard and analytics operations
type DashboardTools struct {
	dashboardUseCase dashboardUC.UseCase
}

// NewDashboardTools creates a new instance of DashboardTools
func NewDashboardTools(dashboardUseCase dashboardUC.UseCase) *DashboardTools {
	return &DashboardTools{
		dashboardUseCase: dashboardUseCase,
	}
}

// GetFinancialSummaryTool implements MCPTool interface for getting comprehensive financial overview
type GetFinancialSummaryTool struct {
	dashboardUseCase dashboardUC.UseCase
}

func NewGetFinancialSummaryTool(dashboardUseCase dashboardUC.UseCase) *GetFinancialSummaryTool {
	return &GetFinancialSummaryTool{
		dashboardUseCase: dashboardUseCase,
	}
}

func (t *GetFinancialSummaryTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_financial_summary",
		Description: "Get comprehensive financial overview for a specific month and year",
		InputSchema: GetFinancialSummarySchema,
	}
}

func (t *GetFinancialSummaryTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	month := int(args["month"].(float64))
	year := int(args["year"].(float64))

	summary, err := t.dashboardUseCase.GetSummary(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial summary: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"summary": map[string]interface{}{
			"month":           month,
			"year":            year,
			"total_income":    summary.TotalIncome,
			"total_expense":   summary.TotalExpense,
			"total_remaining": summary.TotalRemaining,
			"balance_status":  getBalanceStatus(summary.TotalRemaining),
			"generated_at":    time.Now().Format(time.RFC3339),
		},
	}, nil
}

// GetMonthlySummaryTool implements MCPTool interface for getting monthly financial breakdown
type GetMonthlySummaryTool struct {
	dashboardUseCase dashboardUC.UseCase
}

func NewGetMonthlySummaryTool(dashboardUseCase dashboardUC.UseCase) *GetMonthlySummaryTool {
	return &GetMonthlySummaryTool{
		dashboardUseCase: dashboardUseCase,
	}
}

func (t *GetMonthlySummaryTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_monthly_summary",
		Description: "Get monthly financial breakdown with detailed analysis",
		InputSchema: GetMonthlySummarySchema,
	}
}

func (t *GetMonthlySummaryTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	month := int(args["month"].(float64))
	year := int(args["year"].(float64))

	// Get basic summary
	summary, err := t.dashboardUseCase.GetSummary(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly summary: %v", err)
	}

	// Get budget utilization
	budgetUtilization, err := t.dashboardUseCase.SummaryBudgetUsageByMonthYear(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget utilization: %v", err)
	}

	return map[string]interface{}{
		"success": true,
		"monthly_summary": map[string]interface{}{
			"period": map[string]interface{}{
				"month": month,
				"year":  year,
			},
			"financial_overview": map[string]interface{}{
				"total_income":    summary.TotalIncome,
				"total_expense":   summary.TotalExpense,
				"total_remaining": summary.TotalRemaining,
				"savings_rate":    calculateSavingsRate(summary.TotalIncome, summary.TotalRemaining),
			},
			"budget_utilization": budgetUtilization,
			"generated_at":       time.Now().Format(time.RFC3339),
		},
	}, nil
}

// GetBudgetUtilizationTool implements MCPTool interface for getting budget usage analysis
type GetBudgetUtilizationTool struct {
	dashboardUseCase dashboardUC.UseCase
}

func NewGetBudgetUtilizationTool(dashboardUseCase dashboardUC.UseCase) *GetBudgetUtilizationTool {
	return &GetBudgetUtilizationTool{
		dashboardUseCase: dashboardUseCase,
	}
}

func (t *GetBudgetUtilizationTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_budget_utilization",
		Description: "Get budget usage and remaining amounts for a specific period",
		InputSchema: GetBudgetUtilizationSchema,
	}
}

func (t *GetBudgetUtilizationTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	month := int(args["month"].(float64))
	year := int(args["year"].(float64))

	budgetUtilization, err := t.dashboardUseCase.SummaryBudgetUsageByMonthYear(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget utilization: %v", err)
	}

	// Budget utilization data is already processed by the use case
	// Additional analysis could be added here if needed

	return map[string]interface{}{
		"success": true,
		"budget_analysis": map[string]interface{}{
			"period": map[string]interface{}{
				"month": month,
				"year":  year,
			},
			"budgets":      budgetUtilization,
			"generated_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// GetIncomeVsExpensesTool implements MCPTool interface for comparing income against expenses
type GetIncomeVsExpensesTool struct {
	dashboardUseCase dashboardUC.UseCase
}

func NewGetIncomeVsExpensesTool(dashboardUseCase dashboardUC.UseCase) *GetIncomeVsExpensesTool {
	return &GetIncomeVsExpensesTool{
		dashboardUseCase: dashboardUseCase,
	}
}

func (t *GetIncomeVsExpensesTool) GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_income_vs_expenses",
		Description: "Compare income against expenses with trend analysis",
		InputSchema: GetIncomeVsExpensesSchema,
	}
}

func (t *GetIncomeVsExpensesTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	month := int(args["month"].(float64))
	year := int(args["year"].(float64))
	comparisonMonths := 3 // default

	if cm, exists := args["comparison_months"].(float64); exists {
		comparisonMonths = int(cm)
	}

	// Get current month summary
	currentSummary, err := t.dashboardUseCase.GetSummary(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get current month summary: %v", err)
	}

	// Get historical data for comparison
	historicalData := make([]map[string]interface{}, 0, comparisonMonths)
	for i := 1; i <= comparisonMonths; i++ {
		prevMonth := month - i
		prevYear := year

		if prevMonth <= 0 {
			prevMonth += 12
			prevYear--
		}

		prevSummary, err := t.dashboardUseCase.GetSummary(ctx, prevMonth, prevYear)
		if err != nil {
			// Skip months with errors
			continue
		}

		historicalData = append(historicalData, map[string]interface{}{
			"month":           prevMonth,
			"year":            prevYear,
			"total_income":    prevSummary.TotalIncome,
			"total_expense":   prevSummary.TotalExpense,
			"total_remaining": prevSummary.TotalRemaining,
		})
	}

	return map[string]interface{}{
		"success": true,
		"income_vs_expenses": map[string]interface{}{
			"current_period": map[string]interface{}{
				"month":           month,
				"year":            year,
				"total_income":    currentSummary.TotalIncome,
				"total_expense":   currentSummary.TotalExpense,
				"total_remaining": currentSummary.TotalRemaining,
				"savings_rate":    calculateSavingsRate(currentSummary.TotalIncome, currentSummary.TotalRemaining),
			},
			"historical_comparison": historicalData,
			"analysis": map[string]interface{}{
				"balance_status": getBalanceStatus(currentSummary.TotalRemaining),
				"expense_ratio":  calculateExpenseRatio(currentSummary.TotalIncome, currentSummary.TotalExpense),
			},
			"generated_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// GetDashboardTools returns all dashboard-related MCP tools
func (dt *DashboardTools) GetDashboardTools() []interface{} {
	return []interface{}{
		// NewGetFinancialSummaryTool(dt.dashboardUseCase),
		NewGetMonthlySummaryTool(dt.dashboardUseCase),
		NewGetBudgetUtilizationTool(dt.dashboardUseCase),
		NewGetIncomeVsExpensesTool(dt.dashboardUseCase),
	}
}

// Helper functions
func getBalanceStatus(remaining float64) string {
	if remaining > 0 {
		return "positive"
	} else if remaining < 0 {
		return "negative"
	}
	return "neutral"
}

func calculateSavingsRate(income, remaining float64) float64 {
	if income <= 0 {
		return 0.0
	}
	return (remaining / income) * 100
}

func calculateExpenseRatio(income, expense float64) float64 {
	if income <= 0 {
		return 0.0
	}
	return (expense / income) * 100
}
