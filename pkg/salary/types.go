package salary

import "time"

// Entry represents a salary entry that can be either a discount or an addition.
// It supports different calculation types: fixed amounts, percentages, or computed values.
type Entry struct {
	// EntryType indicates whether this is a 'discount' or 'addition' to the salary
	EntryType string

	// Category describes the type of entry (e.g., "INSS", "IRRF", "Extra Hours")
	Category string

	// CalculationType defines how the amount was calculated: 'fixed', 'percentage', or 'computed'
	CalculationType string

	// BaseValue is the base amount used for percentage or computed calculations
	BaseValue float64

	// Multiplier is the factor applied to the base value (e.g., tax rate, hourly multiplier)
	Multiplier float64

	// Amount is the final calculated value for this entry
	Amount float64

	// Description provides additional details about this entry
	Description string

	// ReferenceDate is an optional timestamp for when this entry applies
	ReferenceDate *time.Time
}

// TaxBracket represents a progressive tax bracket with a limit, rate, and optional deduction.
// Used for both INSS and IRRF calculations.
type TaxBracket struct {
	// Limit is the upper bound for this tax bracket (0 means no upper limit)
	Limit float64

	// Rate is the tax percentage applied to the amount within this bracket
	Rate float64

	// Deduction is a fixed amount subtracted from the calculated tax (used in IRRF)
	Deduction float64
}
