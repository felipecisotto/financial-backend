package salary

import (
	"context"
	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto/request"
	"fmt"
	"time"
)

func (uc *useCaseImpl) ActivateTracking(ctx context.Context, incomeID string, req request.ActivateSalaryTrackingRequest) error {
	// 1. Buscar Income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Validar discount mode
	if req.DiscountMode != "CLT" && req.DiscountMode != "custom" {
		return fmt.Errorf("discount_mode must be 'CLT' or 'custom', got '%s'", req.DiscountMode)
	}

	// 3. Validar hourly rate
	if req.HourlyRate <= 0 {
		return fmt.Errorf("hourly_rate must be greater than zero")
	}

	// 4. Criar novo model com os campos atualizados
	updatedIncome, err := models.NewIncome(
		income.ID(),
		income.Description(),
		income.Amount(),
		income.Type(),
		income.DueDay(),
		income.StartDate(),
		income.EndDate(),
		&req.DiscountMode,
		&req.HourlyRate,
	)
	if err != nil {
		return fmt.Errorf("failed to create updated income model: %w", err)
	}

	// 5. Atualizar income
	if err := uc.incomeGateway.Update(ctx, updatedIncome); err != nil {
		return fmt.Errorf("failed to update income: %w", err)
	}

	// 6. Se CLT, criar entries de INSS/IRRF do mês atual
	if req.DiscountMode == "CLT" {
		now := time.Now()
		currentMonth := int(now.Month())
		currentYear := now.Year()

		if err := uc.RecalculateCLT(ctx, incomeID, currentMonth, currentYear); err != nil {
			return fmt.Errorf("failed to create CLT entries: %w", err)
		}
	}

	return nil
}
