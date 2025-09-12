package tools

// GetFinancialSummarySchema defines the JSON schema for getting financial summary
var GetFinancialSummarySchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to get summary for (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to get summary for",
			"minimum": 2000,
			"maximum": 9999
		}
	},
	"required": ["month", "year"]
}`)

// GetMonthlySummarySchema defines the JSON schema for getting monthly breakdown
var GetMonthlySummarySchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to get breakdown for (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to get breakdown for",
			"minimum": 2000,
			"maximum": 9999
		}
	},
	"required": ["month", "year"]
}`)

// GetExpenseAnalyticsSchema defines the JSON schema for expense analysis
var GetExpenseAnalyticsSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to analyze (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to analyze",
			"minimum": 2000,
			"maximum": 9999
		},
		"category": {
			"type": "string",
			"description": "Filter by specific category (optional)"
		}
	},
	"required": ["month", "year"]
}`)

// GetBudgetUtilizationSchema defines the JSON schema for budget usage analysis
var GetBudgetUtilizationSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to get budget utilization for (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to get budget utilization for",
			"minimum": 2000,
			"maximum": 9999
		}
	},
	"required": ["month", "year"]
}`)

// GetIncomeVsExpensesSchema defines the JSON schema for income vs expenses comparison
var GetIncomeVsExpensesSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to compare (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to compare",
			"minimum": 2000,
			"maximum": 9999
		},
		"comparison_months": {
			"type": "integer",
			"description": "Number of previous months to compare against (optional)",
			"minimum": 1,
			"maximum": 12,
			"default": 3
		}
	},
	"required": ["month", "year"]
}`)

// GetCategoryBreakdownSchema defines the JSON schema for category spending analysis
var GetCategoryBreakdownSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"month": {
			"type": "integer",
			"description": "Month to get breakdown for (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Year to get breakdown for",
			"minimum": 2000,
			"maximum": 9999
		},
		"top_categories": {
			"type": "integer",
			"description": "Number of top categories to return (optional)",
			"minimum": 1,
			"maximum": 20,
			"default": 10
		}
	},
	"required": ["month", "year"]
}`)
