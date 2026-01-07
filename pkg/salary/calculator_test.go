package salary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CalculatorTestSuite struct {
	suite.Suite
	calculator Calculator
}

func (s *CalculatorTestSuite) SetupTest() {
	s.calculator = NewCalculator()
}

// TestCalculateProgressiveINSS tests progressive INSS calculation with multiple brackets
func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_FirstBracket() {
	// Test salary in first bracket (up to R$ 1,412.00)
	// Salary: R$ 1,000.00
	// Expected: 1,000.00 * 7.5% = R$ 75.00
	grossAmount := 1000.00
	expected := 75.00

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for first bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_FirstBracketLimit() {
	// Test salary exactly at first bracket limit (R$ 1,412.00)
	// Expected: 1,412.00 * 7.5% = R$ 105.90
	grossAmount := 1412.00
	expected := 105.90

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation at first bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_SecondBracket() {
	// Test salary in second bracket (R$ 1,412.01 - R$ 2,666.68)
	// Salary: R$ 2,000.00
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,000.00 - 1,412.00) * 9% = 588.00 * 9% = 52.92
	// Total: 105.90 + 52.92 = 158.82
	grossAmount := 2000.00
	expected := 158.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for second bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_SecondBracketLimit() {
	// Test salary exactly at second bracket limit (R$ 2,666.68)
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,666.68 - 1,412.00) * 9% = 1,254.68 * 9% = 112.92
	// Total: 105.90 + 112.92 = 218.82
	grossAmount := 2666.68
	expected := 218.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation at second bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_ThirdBracket() {
	// Test salary in third bracket (R$ 2,666.69 - R$ 4,000.03)
	// Salary: R$ 3,500.00
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,666.68 - 1,412.00) * 9% = 1,254.68 * 9% = 112.92
	// Third bracket: (3,500.00 - 2,666.68) * 12% = 833.32 * 12% = 100.00
	// Total: 105.90 + 112.92 + 100.00 = 318.82
	grossAmount := 3500.00
	expected := 318.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for third bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_ThirdBracketLimit() {
	// Test salary exactly at third bracket limit (R$ 4,000.03)
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,666.68 - 1,412.00) * 9% = 1,254.68 * 9% = 112.92
	// Third bracket: (4,000.03 - 2,666.68) * 12% = 1,333.35 * 12% = 160.00
	// Total: 105.90 + 112.92 + 160.00 = 378.82
	grossAmount := 4000.03
	expected := 378.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation at third bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_AboveCeiling() {
	// Test salary above ceiling (R$ 4,000.03+)
	// Salary: R$ 5,000.00
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,666.68 - 1,412.00) * 9% = 1,254.68 * 9% = 112.92
	// Third bracket: (4,000.03 - 2,666.68) * 12% = 1,333.35 * 12% = 160.00
	// Fourth bracket: (5,000.00 - 4,000.03) * 14% = 999.97 * 14% = 140.00
	// Total: 105.90 + 112.92 + 160.00 + 140.00 = 518.82
	grossAmount := 5000.00
	expected := 518.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation above ceiling should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_HighSalary() {
	// Test very high salary
	// Salary: R$ 10,000.00
	// First bracket: 1,412.00 * 7.5% = 105.90
	// Second bracket: (2,666.68 - 1,412.00) * 9% = 1,254.68 * 9% = 112.92
	// Third bracket: (4,000.03 - 2,666.68) * 12% = 1,333.35 * 12% = 160.00
	// Fourth bracket: (10,000.00 - 4,000.03) * 14% = 5,999.97 * 14% = 840.00
	// Total: 105.90 + 112.92 + 160.00 + 840.00 = 1,218.82
	grossAmount := 10000.00
	expected := 1218.82

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for high salary should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_ZeroSalary() {
	// Test zero salary
	grossAmount := 0.0
	expected := 0.0

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for zero salary should be zero")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveINSS_VeryLowSalary() {
	// Test very low salary
	// Salary: R$ 100.00
	// Expected: 100.00 * 7.5% = R$ 7.50
	grossAmount := 100.00
	expected := 7.50

	result := calculateProgressiveINSS(grossAmount)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "INSS calculation for very low salary should be correct")
}

// TestCalculateProgressiveIRRF tests progressive IRRF calculation
func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_TaxExempt() {
	// Test tax-exempt base (up to R$ 2,112.00)
	// Base: R$ 2,000.00
	// Expected: 0.00 (tax exempt)
	taxableBase := 2000.00
	expected := 0.0

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for tax-exempt base should be zero")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_TaxExemptLimit() {
	// Test exactly at tax-exempt limit (R$ 2,112.00)
	// Expected: 0.00 (tax exempt)
	taxableBase := 2112.00
	expected := 0.0

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation at tax-exempt limit should be zero")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_FirstBracket() {
	// Test first tax bracket (R$ 2,112.01 - R$ 2,826.65)
	// Base: R$ 2,500.00
	// Calculation: (2,500.00 * 7.5%) - 158.40 = 187.50 - 158.40 = 29.10
	taxableBase := 2500.00
	expected := 29.10

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for first tax bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_FirstBracketLimit() {
	// Test exactly at first bracket limit (R$ 2,826.65)
	// Calculation: (2,826.65 * 7.5%) - 158.40 = 212.00 - 158.40 = 53.60
	taxableBase := 2826.65
	expected := 53.60

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation at first bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_SecondBracket() {
	// Test second tax bracket (R$ 2,826.66 - R$ 3,751.05)
	// Base: R$ 3,500.00
	// Calculation: (3,500.00 * 15%) - 370.40 = 525.00 - 370.40 = 154.60
	taxableBase := 3500.00
	expected := 154.60

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for second tax bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_SecondBracketLimit() {
	// Test exactly at second bracket limit (R$ 3,751.05)
	// Calculation: (3,751.05 * 15%) - 370.40 = 562.66 - 370.40 = 192.26
	taxableBase := 3751.05
	expected := 192.26

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation at second bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_ThirdBracket() {
	// Test third tax bracket (R$ 3,751.06 - R$ 4,664.68)
	// Base: R$ 4,500.00
	// Calculation: (4,500.00 * 22.5%) - 651.73 = 1,012.50 - 651.73 = 360.77
	taxableBase := 4500.00
	expected := 360.77

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for third tax bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_ThirdBracketLimit() {
	// Test exactly at third bracket limit (R$ 4,664.68)
	// Calculation: (4,664.68 * 22.5%) - 651.73 = 1,049.55 - 651.73 = 397.82
	taxableBase := 4664.68
	expected := 397.82

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation at third bracket limit should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_FourthBracket() {
	// Test fourth tax bracket (R$ 4,664.69+)
	// Base: R$ 6,000.00
	// Calculation: (6,000.00 * 27.5%) - 884.96 = 1,650.00 - 884.96 = 765.04
	taxableBase := 6000.00
	expected := 765.04

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for fourth tax bracket should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_HighBase() {
	// Test very high base
	// Base: R$ 10,000.00
	// Calculation: (10,000.00 * 27.5%) - 884.96 = 2,750.00 - 884.96 = 1,865.04
	taxableBase := 10000.00
	expected := 1865.04

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for high base should be correct")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_ZeroBase() {
	// Test zero base
	taxableBase := 0.0
	expected := 0.0

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation for zero base should be zero")
}

func (s *CalculatorTestSuite) TestCalculateProgressiveIRRF_NegativeResultHandling() {
	// Test case where deduction is larger than tax (should return 0)
	// Base: R$ 2,113.00 (just above exempt limit)
	// Calculation: (2,113.00 * 7.5%) - 158.40 = 158.475 - 158.40 = 0.075 = 0.07 (rounded)
	taxableBase := 2113.00
	expected := 0.07

	result := calculateProgressiveIRRF(taxableBase)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "IRRF calculation should handle small positive results")
}

// TestCalculateExtraHours tests extra hours calculation
func (s *CalculatorTestSuite) TestCalculateExtraHours_SimpleCalculation() {
	// Simple calculation: 30 hours * R$ 10.00 = R$ 300.00
	hourlyRate := 10.00
	hours := 30.0
	expected := 300.00

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation should be correct")
}

func (s *CalculatorTestSuite) TestCalculateExtraHours_DifferentRates() {
	// Test with different hourly rate
	// 20 hours * R$ 25.50 = R$ 510.00
	hourlyRate := 25.50
	hours := 20.0
	expected := 510.00

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation with different rate should be correct")
}

func (s *CalculatorTestSuite) TestCalculateExtraHours_DecimalHours() {
	// Test with decimal hours
	// 15.5 hours * R$ 30.00 = R$ 465.00
	hourlyRate := 30.00
	hours := 15.5
	expected := 465.00

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation with decimal hours should be correct")
}

func (s *CalculatorTestSuite) TestCalculateExtraHours_ZeroHours() {
	// Test with zero hours
	hourlyRate := 20.00
	hours := 0.0
	expected := 0.0

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation with zero hours should be zero")
}

func (s *CalculatorTestSuite) TestCalculateExtraHours_HighRate() {
	// Test with high hourly rate
	// 10 hours * R$ 150.00 = R$ 1,500.00
	hourlyRate := 150.00
	hours := 10.0
	expected := 1500.00

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation with high rate should be correct")
}

func (s *CalculatorTestSuite) TestCalculateExtraHours_SmallHours() {
	// Test with very small hours
	// 0.5 hours * R$ 20.00 = R$ 10.00
	hourlyRate := 20.00
	hours := 0.5
	expected := 10.00

	result := s.calculator.CalculateExtraHours(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "Extra hours calculation with small hours should be correct")
}

// TestCalculateOnCall tests on-call hours calculation
func (s *CalculatorTestSuite) TestCalculateOnCall_SimpleCalculation() {
	// 1/3 calculation: (R$ 30.00 / 3) * 24 hours = R$ 10.00 * 24 = R$ 240.00
	hourlyRate := 30.00
	hours := 24.0
	expected := 240.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation should be correct")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_DifferentRates() {
	// Test with different rate
	// (R$ 45.00 / 3) * 18 hours = R$ 15.00 * 18 = R$ 270.00
	hourlyRate := 45.00
	hours := 18.0
	expected := 270.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation with different rate should be correct")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_DecimalHours() {
	// Test with decimal hours
	// (R$ 60.00 / 3) * 12.5 hours = R$ 20.00 * 12.5 = R$ 250.00
	hourlyRate := 60.00
	hours := 12.5
	expected := 250.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation with decimal hours should be correct")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_ZeroHours() {
	// Test with zero hours
	hourlyRate := 30.00
	hours := 0.0
	expected := 0.0

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation with zero hours should be zero")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_HighRate() {
	// Test with high hourly rate
	// (R$ 150.00 / 3) * 20 hours = R$ 50.00 * 20 = R$ 1,000.00
	hourlyRate := 150.00
	hours := 20.0
	expected := 1000.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation with high rate should be correct")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_SmallHours() {
	// Test with very small hours
	// (R$ 30.00 / 3) * 1 hour = R$ 10.00 * 1 = R$ 10.00
	hourlyRate := 30.00
	hours := 1.0
	expected := 10.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation with small hours should be correct")
}

func (s *CalculatorTestSuite) TestCalculateOnCall_PrecisionCheck() {
	// Test precision with rate that doesn't divide evenly by 3
	// (R$ 100.00 / 3) * 3 hours = R$ 33.333... * 3 = R$ 100.00
	hourlyRate := 100.00
	hours := 3.0
	expected := 100.00

	result := s.calculator.CalculateOnCall(hourlyRate, hours)

	assert.Equal(s.T(), expected, roundToTwoDecimals(result), "On-call calculation should handle precision correctly")
}

// TestCalculateCLTDiscounts - Integration test for full CLT discount calculation
func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_LowSalary() {
	// Test with low salary that results in both INSS and IRRF
	// Gross: R$ 3,000.00
	// INSS: (1,412.00 * 7.5%) + (1,254.68 * 9%) + (333.32 * 12%) = 105.90 + 112.92 + 40.00 = 258.82
	// IRRF Base: 3,000.00 - 258.82 = 2,741.18
	// IRRF: (2,741.18 * 7.5%) - 158.40 = 205.59 - 158.40 = 47.19
	grossAmount := 3000.00
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2, "Should return both INSS and IRRF entries")

	// Check INSS entry
	assert.Equal(s.T(), "discount", entries[0].EntryType)
	assert.Equal(s.T(), "INSS", entries[0].Category)
	assert.Equal(s.T(), "computed", entries[0].CalculationType)
	assert.Equal(s.T(), grossAmount, entries[0].BaseValue)
	assert.Equal(s.T(), 258.82, entries[0].Amount)

	// Check IRRF entry
	assert.Equal(s.T(), "discount", entries[1].EntryType)
	assert.Equal(s.T(), "IRRF", entries[1].Category)
	assert.Equal(s.T(), "computed", entries[1].CalculationType)
	assert.InDelta(s.T(), 2741.18, entries[1].BaseValue, 0.01, "IRRF base value should be approximately correct")
	assert.Equal(s.T(), 47.19, entries[1].Amount)
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_HighSalary() {
	// Test with high salary
	// Gross: R$ 8,000.00
	// INSS: (1,412.00 * 7.5%) + (1,254.68 * 9%) + (1,333.35 * 12%) + (3,999.97 * 14%) = 105.90 + 112.92 + 160.00 + 560.00 = 938.82
	// IRRF Base: 8,000.00 - 938.82 = 7,061.18
	// IRRF: (7,061.18 * 27.5%) - 884.96 = 1,941.82 - 884.96 = 1,056.86
	grossAmount := 8000.00
	month := 6
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2, "Should return both INSS and IRRF entries")

	// Check INSS entry
	assert.Equal(s.T(), "INSS", entries[0].Category)
	assert.Equal(s.T(), 938.82, entries[0].Amount)

	// Check IRRF entry
	assert.Equal(s.T(), "IRRF", entries[1].Category)
	assert.Equal(s.T(), 1056.86, entries[1].Amount)
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_TaxExempt() {
	// Test with salary that results in INSS but no IRRF
	// Gross: R$ 2,000.00
	// INSS: (1,412.00 * 7.5%) + (588.00 * 9%) = 105.90 + 52.92 = 158.82
	// IRRF Base: 2,000.00 - 158.82 = 1,841.18 (below exempt limit of 2,112.00)
	// IRRF: 0.00
	grossAmount := 2000.00
	month := 3
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 1, "Should return only INSS entry (IRRF is zero)")

	// Check INSS entry
	assert.Equal(s.T(), "INSS", entries[0].Category)
	assert.Equal(s.T(), 158.82, entries[0].Amount)
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_ZeroSalary() {
	// Test with zero salary
	grossAmount := 0.0
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 0, "Should return empty array for zero salary")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_NegativeSalary() {
	// Test error handling for negative salary
	grossAmount := -1000.00
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), entries)
	assert.Contains(s.T(), err.Error(), "gross amount cannot be negative")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_InvalidMonth() {
	// Test error handling for invalid month
	grossAmount := 3000.00
	month := 13
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), entries)
	assert.Contains(s.T(), err.Error(), "month must be between 1 and 12")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_InvalidMonthZero() {
	// Test error handling for zero month
	grossAmount := 3000.00
	month := 0
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), entries)
	assert.Contains(s.T(), err.Error(), "month must be between 1 and 12")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_InvalidYear() {
	// Test error handling for invalid year
	grossAmount := 3000.00
	month := 1
	year := 1800

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), entries)
	assert.Contains(s.T(), err.Error(), "year must be between 1900 and 2100")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_InvalidYearFuture() {
	// Test error handling for future year
	grossAmount := 3000.00
	month := 1
	year := 2200

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), entries)
	assert.Contains(s.T(), err.Error(), "year must be between 1900 and 2100")
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_EdgeCaseAtINSSCeiling() {
	// Test exactly at INSS ceiling (R$ 4,000.03)
	// INSS: (1,412.00 * 7.5%) + (1,254.68 * 9%) + (1,333.35 * 12%) = 105.90 + 112.92 + 160.00 = 378.82
	// IRRF Base: 4,000.03 - 378.82 = 3,621.21
	// IRRF: (3,621.21 * 15%) - 370.40 = 543.18 - 370.40 = 172.78
	grossAmount := 4000.03
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2, "Should return both INSS and IRRF entries")

	// Check INSS entry
	assert.Equal(s.T(), "INSS", entries[0].Category)
	assert.Equal(s.T(), 378.82, entries[0].Amount)

	// Check IRRF entry
	assert.Equal(s.T(), "IRRF", entries[1].Category)
	assert.Equal(s.T(), 172.78, entries[1].Amount)
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_MinimumWage() {
	// Test with minimum wage (R$ 1,412.00)
	// INSS: 1,412.00 * 7.5% = 105.90
	// IRRF Base: 1,412.00 - 105.90 = 1,306.10 (below exempt limit)
	// IRRF: 0.00
	grossAmount := 1412.00
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 1, "Should return only INSS entry")

	// Check INSS entry
	assert.Equal(s.T(), "INSS", entries[0].Category)
	assert.Equal(s.T(), 105.90, entries[0].Amount)
}

func (s *CalculatorTestSuite) TestCalculateCLTDiscounts_EntryStructure() {
	// Test the complete structure of returned entries
	grossAmount := 5000.00
	month := 1
	year := 2024

	entries, err := s.calculator.CalculateCLTDiscounts(grossAmount, month, year)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2)

	// Check INSS entry structure
	inss := entries[0]
	assert.Equal(s.T(), "discount", inss.EntryType)
	assert.Equal(s.T(), "INSS", inss.Category)
	assert.Equal(s.T(), "computed", inss.CalculationType)
	assert.Greater(s.T(), inss.BaseValue, 0.0)
	assert.Equal(s.T(), float64(0), inss.Multiplier) // Progressive, no single multiplier
	assert.Greater(s.T(), inss.Amount, 0.0)
	assert.Contains(s.T(), inss.Description, "INSS")
	assert.Nil(s.T(), inss.ReferenceDate)

	// Check IRRF entry structure
	irrf := entries[1]
	assert.Equal(s.T(), "discount", irrf.EntryType)
	assert.Equal(s.T(), "IRRF", irrf.Category)
	assert.Equal(s.T(), "computed", irrf.CalculationType)
	assert.Greater(s.T(), irrf.BaseValue, 0.0)
	assert.Equal(s.T(), float64(0), irrf.Multiplier) // Progressive, no single multiplier
	assert.Greater(s.T(), irrf.Amount, 0.0)
	assert.Contains(s.T(), irrf.Description, "IRRF")
	assert.Nil(s.T(), irrf.ReferenceDate)
}

// TestRoundToTwoDecimals tests the rounding utility function
func (s *CalculatorTestSuite) TestRoundToTwoDecimals_RoundUp() {
	value := 123.456
	expected := 123.46

	result := roundToTwoDecimals(value)

	assert.Equal(s.T(), expected, result, "Should round up to 2 decimal places")
}

func (s *CalculatorTestSuite) TestRoundToTwoDecimals_RoundDown() {
	value := 123.454
	expected := 123.45

	result := roundToTwoDecimals(value)

	assert.Equal(s.T(), expected, result, "Should round down to 2 decimal places")
}

func (s *CalculatorTestSuite) TestRoundToTwoDecimals_ExactValue() {
	value := 123.45
	expected := 123.45

	result := roundToTwoDecimals(value)

	assert.Equal(s.T(), expected, result, "Should keep exact 2 decimal places")
}

func (s *CalculatorTestSuite) TestRoundToTwoDecimals_ZeroDecimals() {
	value := 123.0
	expected := 123.0

	result := roundToTwoDecimals(value)

	assert.Equal(s.T(), expected, result, "Should handle values with no decimals")
}

func (s *CalculatorTestSuite) TestRoundToTwoDecimals_NegativeValue() {
	value := -123.456
	expected := -123.46

	result := roundToTwoDecimals(value)

	assert.Equal(s.T(), expected, result, "Should round negative values correctly")
}

// Run the test suite
func TestCalculatorTestSuite(t *testing.T) {
	suite.Run(t, new(CalculatorTestSuite))
}
