package tools

// CreateExpenseSchema defines the JSON schema for creating expenses
var CreateExpenseSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Description of the expense",
			"minLength": 1
		},
		"amount": {
			"type": "number",
			"description": "Amount of the expense",
			"minimum": 0.01
		},
		"type": {
			"type": "string",
			"description": "Type of expense (e.g., DEBIT, CREDIT)",
			"enum": ["DEBIT", "CREDIT"]
		},
		"budget_id": {
			"type": "string",
			"description": "ID of the budget this expense belongs to (optional)"
		},
		"recurrency": {
			"type": "string",
			"description": "Recurrence pattern (optional)",
			"enum": ["MONTHLY", "WEEKLY", "YEARLY"]
		},
		"method": {
			"type": "string",
			"description": "Payment method",
			"enum": ["CASH", "CREDIT_CARD", "DEBIT_CARD", "PIX", "BANK_TRANSFER"]
		},
		"installments": {
			"type": "integer",
			"description": "Number of installments (optional)",
			"minimum": 1
		},
		"due_day": {
			"type": "integer",
			"description": "Day of the month when expense is due",
			"minimum": 1,
			"maximum": 31
		},
		"start_date": {
			"type": "string",
			"format": "date",
			"description": "Start date of the expense"
		},
		"end_date": {
			"type": "string",
			"format": "date",
			"description": "End date of the expense (optional)"
		}
	},
	"required": ["description", "amount", "type", "method", "due_day", "start_date"]
}`)

// GetExpensesSchema defines the JSON schema for listing expenses
var GetExpensesSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Filter by description (partial match)"
		},
		"type": {
			"type": "string",
			"description": "Filter by expense type",
			"enum": ["DEBIT", "CREDIT"]
		},
		"category": {
			"type": "string",
			"description": "Filter by category"
		},
		"budget_id": {
			"type": "string",
			"description": "Filter by budget ID"
		},
		"recurrency": {
			"type": "string",
			"description": "Filter by recurrence pattern",
			"enum": ["MONTHLY", "WEEKLY", "YEARLY"]
		},
		"method": {
			"type": "string",
			"description": "Filter by payment method",
			"enum": ["CASH", "CREDIT_CARD", "DEBIT_CARD", "PIX", "BANK_TRANSFER"]
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

// GetExpenseByIdSchema defines the JSON schema for getting a specific expense
var GetExpenseByIdSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the expense",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)

// UpdateExpenseSchema defines the JSON schema for updating expenses
var UpdateExpenseSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the expense to update",
			"minLength": 1
		},
		"description": {
			"type": "string",
			"description": "New description of the expense"
		},
		"amount": {
			"type": "number",
			"description": "New amount of the expense",
			"minimum": 0.01
		},
		"category": {
			"type": "string",
			"description": "New category of the expense"
		},
		"start_date": {
			"type": "string",
			"format": "date",
			"description": "New start date of the expense"
		},
		"end_date": {
			"type": "string",
			"format": "date",
			"description": "New end date of the expense"
		},
		"installments": {
			"type": "integer",
			"description": "New number of installments",
			"minimum": 1
		},
		"due_date": {
			"type": "string",
			"format": "date",
			"description": "New due date of the expense"
		},
		"statement_day": {
			"type": "integer",
			"description": "New statement day",
			"minimum": 1,
			"maximum": 31
		}
	},
	"required": ["id"]
}`)

// DeleteExpenseSchema defines the JSON schema for deleting expenses
var DeleteExpenseSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the expense to delete",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)
