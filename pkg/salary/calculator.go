package salary

// Calculator defines the interface for salary-related calculations.
// It provides methods for calculating CLT discounts (INSS and IRRF),
// extra hours compensation, and on-call pay.
type Calculator interface {
	// CalculateCLTDiscounts calculates INSS and IRRF deductions for a given gross salary.
	// It returns a slice of Entry objects representing each deduction.
	// Parameters:
	//   - grossAmount: The gross salary amount
	//   - month: The month for the calculation (1-12)
	//   - year: The year for the calculation
	CalculateCLTDiscounts(grossAmount float64, month, year int) ([]Entry, error)

	// CalculateExtraHours calculates the compensation for extra hours worked.
	// Formula: hourlyRate * hours
	// Parameters:
	//   - hourlyRate: The hourly rate for the employee
	//   - hours: Number of extra hours worked
	CalculateExtraHours(hourlyRate float64, hours float64) float64

	// CalculateOnCall calculates the compensation for on-call hours.
	// Formula: (hourlyRate / 3) * hours
	// Parameters:
	//   - hourlyRate: The hourly rate for the employee
	//   - hours: Number of on-call hours
	CalculateOnCall(hourlyRate float64, hours float64) float64
}

// calculator is the concrete implementation of the Calculator interface.
type calculator struct{}

// NewCalculator creates and returns a new Calculator instance.
func NewCalculator() Calculator {
	return &calculator{}
}

// CalculateExtraHours calculates extra hours compensation.
// Extra hours are paid at 100% of the hourly rate.
func (c *calculator) CalculateExtraHours(hourlyRate float64, hours float64) float64 {
	return hourlyRate * hours
}

// CalculateOnCall calculates on-call compensation.
// On-call hours are paid at 33.33% (1/3) of the hourly rate.
func (c *calculator) CalculateOnCall(hourlyRate float64, hours float64) float64 {
	return (hourlyRate / 3.0) * hours
}
