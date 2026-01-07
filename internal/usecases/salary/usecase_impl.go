package salary

import (
	"financial-backend/internal/gateways"
	"financial-backend/pkg/salary"
)

type useCaseImpl struct {
	incomeGateway      gateways.IncomeGateway
	salaryEntryGateway gateways.SalaryEntryGateway
	calculator         salary.Calculator
}

// NewSalaryUseCase creates a new instance of the salary use case
func NewSalaryUseCase(
	incomeGW gateways.IncomeGateway,
	salaryEntryGW gateways.SalaryEntryGateway,
	calc salary.Calculator,
) UseCase {
	return &useCaseImpl{
		incomeGateway:      incomeGW,
		salaryEntryGateway: salaryEntryGW,
		calculator:         calc,
	}
}
