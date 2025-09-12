package tools

// CreateIncomeSchema defines the JSON schema for creating income records
var CreateIncomeSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Description of the income",
			"minLength": 1
		},
		"amount": {
			"type": "number",
			"description": "Amount of the income",
			"minimum": 0.01
		},
		"type": {
			"type": "string",
			"description": "Type of income",
			"enum": ["FIXED", "VARIABLE", "BONUS", "OTHER"]
		},
		"due_day": {
			"type": "integer",
			"description": "Day of the month when income is received",
			"minimum": 1,
			"maximum": 31
		},
		"start_date": {
			"type": "string",
			"format": "date",
			"description": "Start date of the income"
		},
		"end_date": {
			"type": "string",
			"format": "date",
			"description": "End date of the income (optional)"
		}
	},
	"required": ["description", "amount", "type", "due_day", "start_date"]
}`)

// GetIncomesSchema defines the JSON schema for listing income records
var GetIncomesSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"description": {
			"type": "string",
			"description": "Filter by description (partial match)"
		},
		"type": {
			"type": "string",
			"description": "Filter by income type",
			"enum": ["FIXED", "VARIABLE", "BONUS", "OTHER"]
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

// GetIncomeByIdSchema defines the JSON schema for getting a specific income record
var GetIncomeByIdSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the income record",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)

// UpdateIncomeSchema defines the JSON schema for updating income records
var UpdateIncomeSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the income to update",
			"minLength": 1
		},
		"description": {
			"type": "string",
			"description": "New description of the income"
		},
		"amount": {
			"type": "number",
			"description": "New amount of the income",
			"minimum": 0.01
		},
		"type": {
			"type": "string",
			"description": "New type of income",
			"enum": ["FIXED", "VARIABLE", "BONUS", "OTHER"]
		},
		"category": {
			"type": "string",
			"description": "New category of the income"
		},
		"date": {
			"type": "string",
			"format": "date",
			"description": "New date of the income"
		}
	},
	"required": ["id"]
}`)

// DeleteIncomeSchema defines the JSON schema for deleting income records
var DeleteIncomeSchema = mustParseSchema(`{
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "Unique identifier of the income to delete",
			"minLength": 1
		}
	},
	"required": ["id"]
}`)
