package tools

// CreateBudgetSchema defines the JSON schema for creating budgets
// CreateBudgetSchema defines the JSON schema for creating budgets
var CreateBudgetSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Description of the budget",
			"minLength": 1
		},
		"amount": {
			"type": "number",
			"description": "Amount allocated for the budget",
			"minimum": 0.01
		},
		"end_date": {
			"type": "string",
			"format": "date",
			"description": "End date of the budget (optional)"
		}
	},
	"required": ["description", "amount"]
}`)

// GetBudgetsSchema defines the JSON schema for listing budgets
var GetBudgetsSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Filter by description (partial match)"
		},
		"status": {
			"type": "string",
			"description": "Filter by budget status",
			"enum": ["ACTIVE", "INACTIVE", "COMPLETED"]
		},
		"page": {
			"type": "integer",
			"description": "Page number for pagination",
			"minimum": 1,
			"default": 1
		},
		"limit": {
			"type": "integer",
			"description": "Number of items per page",
			"minimum": 1,
			"maximum": 100,
			"default": 10
		}
	}
}`)

// GetBudgetByIdSchema defines the JSON schema for getting a specific budget
var GetBudgetByIdSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the budget",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)

// UpdateBudgetSchema defines the JSON schema for updating budgets
var UpdateBudgetSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the budget to update",
			"minLength": 1
		},
		"end_date": {
			"type": "string",
			"format": "date",
			"description": "New end date for the budget"
		}
	},
	"required": ["id", "end_date"]
}`)

// DeleteBudgetSchema defines the JSON schema for deleting budgets
var DeleteBudgetSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the budget to delete",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)

// GetBudgetMovementsSchema defines the JSON schema for retrieving budget movements
var GetBudgetMovementsSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"budget_id": {
			"type": "string",
			"description": "Filter by budget ID"
		},
		"movement_type": {
			"type": "string",
			"description": "Filter by movement type",
			"enum": ["ALLOCATION", "EXPENSE", "TRANSFER"]
		},
		"origin": {
			"type": "string",
			"description": "Filter by origin of the movement"
		},
		"month": {
			"type": "integer",
			"description": "Filter by month (1-12)",
			"minimum": 1,
			"maximum": 12
		},
		"year": {
			"type": "integer",
			"description": "Filter by year",
			"minimum": 2000,
			"maximum": 9999
		},
		"page": {
			"type": "integer",
			"description": "Page number for pagination",
			"minimum": 1,
			"default": 1
		},
		"limit": {
			"type": "integer",
			"description": "Number of items per page",
			"minimum": 1,
			"maximum": 100,
			"default": 10
		}
	}
}`)
