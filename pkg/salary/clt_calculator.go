package salary

import (
	"fmt"
	"math"
)

// INSS tax brackets for 2024/2025 (progressive calculation)
var inssBrackets = []TaxBracket{
	{Limit: 1412.00, Rate: 0.075, Deduction: 0},
	{Limit: 2666.68, Rate: 0.09, Deduction: 0},
	{Limit: 4000.03, Rate: 0.12, Deduction: 0},
	{Limit: 0, Rate: 0.14, Deduction: 0}, // No upper limit for the last bracket
}

// IRRF tax brackets for 2024/2025 (progressive calculation with deductions)
var irrfBrackets = []TaxBracket{
	{Limit: 2112.00, Rate: 0.0, Deduction: 0.0},      // Exempt
	{Limit: 2826.65, Rate: 0.075, Deduction: 158.40}, // 7.5%
	{Limit: 3751.05, Rate: 0.15, Deduction: 370.40},  // 15%
	{Limit: 4664.68, Rate: 0.225, Deduction: 651.73}, // 22.5%
	{Limit: 0, Rate: 0.275, Deduction: 884.96},       // 27.5% - No upper limit
}

// CalculateCLTDiscounts calculates both INSS and IRRF deductions using progressive tax brackets.
// The IRRF calculation uses the gross amount minus INSS as the taxable base.
func (c *calculator) CalculateCLTDiscounts(grossAmount float64, month, year int) ([]Entry, error) {
	if grossAmount < 0 {
		return nil, fmt.Errorf("gross amount cannot be negative")
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("month must be between 1 and 12")
	}

	if year < 1900 || year > 2100 {
		return nil, fmt.Errorf("year must be between 1900 and 2100")
	}

	entries := make([]Entry, 0, 2)

	// Calculate INSS (progressive)
	inssAmount := calculateProgressiveINSS(grossAmount)
	if inssAmount > 0 {
		entries = append(entries, Entry{
			EntryType:       "discount",
			Category:        "INSS",
			CalculationType: "computed",
			BaseValue:       grossAmount,
			Multiplier:      0, // Progressive, no single multiplier
			Amount:          roundToTwoDecimals(inssAmount),
			Description:     "INSS - Instituto Nacional do Seguro Social (progressive)",
			ReferenceDate:   nil,
		})
	}

	// Calculate IRRF base (gross - INSS)
	irrfBase := grossAmount - inssAmount

	// Calculate IRRF (progressive with deductions)
	irrfAmount := calculateProgressiveIRRF(irrfBase)
	if irrfAmount > 0 {
		entries = append(entries, Entry{
			EntryType:       "discount",
			Category:        "IRRF",
			CalculationType: "computed",
			BaseValue:       irrfBase,
			Multiplier:      0, // Progressive, no single multiplier
			Amount:          roundToTwoDecimals(irrfAmount),
			Description:     "IRRF - Imposto de Renda Retido na Fonte (progressive)",
			ReferenceDate:   nil,
		})
	}

	return entries, nil
}

// calculateProgressiveINSS calculates INSS using progressive tax brackets.
// Each portion of the salary is taxed at its corresponding rate.
func calculateProgressiveINSS(grossAmount float64) float64 {
	totalINSS := 0.0
	previousLimit := 0.0

	for _, bracket := range inssBrackets {
		// For the last bracket (Limit = 0), use the remaining amount
		if bracket.Limit == 0 {
			if grossAmount > previousLimit {
				portionInBracket := grossAmount - previousLimit
				totalINSS += portionInBracket * bracket.Rate
			}
			break
		}

		// If salary is below this bracket's limit, calculate only the portion within it
		if grossAmount <= bracket.Limit {
			portionInBracket := grossAmount - previousLimit
			totalINSS += portionInBracket * bracket.Rate
			break
		}

		// Salary exceeds this bracket, calculate the full bracket amount
		portionInBracket := bracket.Limit - previousLimit
		totalINSS += portionInBracket * bracket.Rate
		previousLimit = bracket.Limit
	}

	return totalINSS
}

// calculateProgressiveIRRF calculates IRRF using progressive tax brackets and deductions.
// The calculation applies the effective rate and deduction for the appropriate bracket.
func calculateProgressiveIRRF(taxableBase float64) float64 {
	// Find the appropriate bracket
	for _, bracket := range irrfBrackets {
		// Check if amount falls within this bracket
		if bracket.Limit == 0 || taxableBase <= bracket.Limit {
			// Apply the bracket's rate and deduction
			irrfAmount := (taxableBase * bracket.Rate) - bracket.Deduction
			if irrfAmount < 0 {
				return 0
			}
			return irrfAmount
		}
	}

	return 0
}

// roundToTwoDecimals rounds a float64 value to 2 decimal places.
func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
