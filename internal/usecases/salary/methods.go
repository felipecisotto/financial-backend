package salary

import (
	"context"
	"errors"
	"financial-backend/internal/domain/models"
	"financial-backend/internal/dto/request"
	"financial-backend/internal/dto/response"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// parseReferenceDate converts a string date in ISO format to *time.Time
func parseReferenceDate(dateStr *string) (*time.Time, error) {
	if dateStr == nil {
		return nil, nil
	}

	parsedTime, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid reference_date format, expected YYYY-MM-DD: %w", err)
	}

	return &parsedTime, nil
}

// GetSummary retrieves the monthly salary summary for an income
func (uc *useCaseImpl) GetSummary(ctx context.Context, incomeID string, month, year int) (*response.SalarySummaryResponse, error) {
	// 1. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Buscar salary entries para o período
	entries, err := uc.salaryEntryGateway.FindByIncomeAndPeriod(ctx, incomeID, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get salary entries: %w", err)
	}

	// 3. Inicializar arrays vazios se nil
	if entries == nil {
		entries = []models.SalaryEntry{}
	}

	// 4. Converter entries para response
	entryResponses := make([]response.SalaryEntryResponse, 0, len(entries))
	var totalAdditions, totalDiscounts float64

	for _, entry := range entries {
		entryResp := response.SalaryEntryResponse{
			ID:              entry.ID(),
			IncomeID:        entry.IncomeID(),
			Month:           entry.Month(),
			Year:            entry.Year(),
			EntryType:       entry.EntryType(),
			Category:        entry.Category(),
			CalculationType: entry.CalculationType(),
			BaseValue:       entry.BaseValue(),
			Multiplier:      entry.Multiplier(),
			Amount:          entry.Amount(),
			Description:     entry.Description(),
			ReferenceDate:   entry.ReferenceDate(),
			CreatedAt:       entry.CreatedAt(),
			UpdatedAt:       entry.UpdatedAt(),
		}
		entryResponses = append(entryResponses, entryResp)

		// Calcular totais
		if entry.IsAddition() {
			totalAdditions += entry.Amount()
		} else if entry.IsDiscount() {
			totalDiscounts += entry.Amount()
		}
	}

	// 5. Calcular net amount
	netAmount := income.Amount() + totalAdditions - totalDiscounts

	// 6. Montar response (sempre retornar estrutura válida)
	summary := &response.SalarySummaryResponse{
		IncomeID:       incomeID,
		Description:    income.Description(),
		Month:          month,
		Year:           year,
		BaseAmount:     income.Amount(),
		Entries:        entryResponses,
		TotalAdditions: totalAdditions,
		TotalDiscounts: totalDiscounts,
		NetAmount:      netAmount,
	}

	return summary, nil
}

// AddEntry adds a generic salary entry
func (uc *useCaseImpl) AddEntry(ctx context.Context, incomeID string, req request.CreateSalaryEntryRequest) (*response.SalarySummaryResponse, error) {
	// 1. Validar que o income existe
	_, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Calcular amount
	amount := req.BaseValue * req.Multiplier

	// 3. Parse reference date
	referenceDate, err := parseReferenceDate(req.ReferenceDate)
	if err != nil {
		return nil, err
	}

	// 4. Criar salary entry model
	entryID := uuid.New().String()
	entry, err := models.NewSalaryEntry(
		entryID,
		incomeID,
		req.Month,
		req.Year,
		req.EntryType,
		req.Category,
		req.CalculationType,
		req.BaseValue,
		req.Multiplier,
		amount,
		req.Description,
		referenceDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create salary entry model: %w", err)
	}

	// 5. Salvar entry
	if err := uc.salaryEntryGateway.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to save salary entry: %w", err)
	}

	// 6. Retornar summary atualizado
	return uc.GetSummary(ctx, incomeID, req.Month, req.Year)
}

// AddExtraHours adds extra hours to the salary
func (uc *useCaseImpl) AddExtraHours(ctx context.Context, incomeID string, req request.AddExtraHoursRequest) (*response.SalarySummaryResponse, error) {
	// 1. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Validar que tem salary tracking
	if !income.HasSalaryTracking() {
		return nil, errors.New("income does not have salary tracking enabled")
	}

	// 3. Validar hourly rate
	if income.HourlyRate() == nil {
		return nil, errors.New("income does not have hourly rate configured")
	}

	hourlyRate := *income.HourlyRate()

	// 4. Se rate não fornecido, usar hourly rate
	if req.Rate == nil {
		req.Rate = &hourlyRate
	}

	// 5. Calcular valor das horas extras
	amount := uc.calculator.CalculateExtraHours(*req.Rate, req.Hours)

	// 6. Criar entry
	entryID := uuid.New().String()
	entry, err := models.NewSalaryEntry(
		entryID,
		incomeID,
		req.Month,
		req.Year,
		"addition",
		"Extra Hours",
		"computed",
		*req.Rate,
		req.Hours,
		amount,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create extra hours entry: %w", err)
	}

	// 7. Salvar entry
	if err := uc.salaryEntryGateway.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to save extra hours entry: %w", err)
	}

	// 8. Se CLT, recalcular descontos
	if income.DiscountMode() != nil && *income.DiscountMode() == "CLT" {
		if err := uc.RecalculateCLT(ctx, incomeID, req.Month, req.Year); err != nil {
			return nil, fmt.Errorf("failed to recalculate CLT: %w", err)
		}
	}

	// 9. Retornar summary atualizado
	return uc.GetSummary(ctx, incomeID, req.Month, req.Year)
}

// AddOnCall adds on-call hours to the salary
func (uc *useCaseImpl) AddOnCall(ctx context.Context, incomeID string, req request.AddOnCallRequest) (*response.SalarySummaryResponse, error) {
	// 1. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Validar que tem salary tracking
	if !income.HasSalaryTracking() {
		return nil, errors.New("income does not have salary tracking enabled")
	}

	// 3. Validar hourly rate
	if income.HourlyRate() == nil {
		return nil, errors.New("income does not have hourly rate configured")
	}

	hourlyRate := *income.HourlyRate()

	// 4. Calcular valor do on-call
	amount := uc.calculator.CalculateOnCall(hourlyRate, req.Hours)

	// 5. Criar entry
	entryID := uuid.New().String()
	entry, err := models.NewSalaryEntry(
		entryID,
		incomeID,
		req.Month,
		req.Year,
		"addition",
		"On-Call",
		"computed",
		hourlyRate/3.0, // base value é o hourly rate / 3
		req.Hours,
		amount,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create on-call entry: %w", err)
	}

	// 6. Salvar entry
	if err := uc.salaryEntryGateway.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to save on-call entry: %w", err)
	}

	// 7. Se CLT, recalcular descontos
	if income.DiscountMode() != nil && *income.DiscountMode() == "CLT" {
		if err := uc.RecalculateCLT(ctx, incomeID, req.Month, req.Year); err != nil {
			return nil, fmt.Errorf("failed to recalculate CLT: %w", err)
		}
	}

	// 8. Retornar summary atualizado
	return uc.GetSummary(ctx, incomeID, req.Month, req.Year)
}

// RemoveEntry removes a salary entry by ID
func (uc *useCaseImpl) RemoveEntry(ctx context.Context, entryID string) error {
	// 1. Buscar entry para pegar incomeID, month e year
	entry, err := uc.salaryEntryGateway.FindByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("failed to find entry: %w", err)
	}

	incomeID := entry.IncomeID()
	month := entry.Month()
	year := entry.Year()

	// 2. Deletar entry
	if err := uc.salaryEntryGateway.Delete(ctx, entryID); err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	// 3. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return fmt.Errorf("failed to get income: %w", err)
	}

	// 4. Se CLT, recalcular descontos
	if income.DiscountMode() != nil && *income.DiscountMode() == "CLT" {
		if err := uc.RecalculateCLT(ctx, incomeID, month, year); err != nil {
			return fmt.Errorf("failed to recalculate CLT: %w", err)
		}
	}

	return nil
}

// GetHistory retrieves the annual salary history
func (uc *useCaseImpl) GetHistory(ctx context.Context, incomeID string, year int) (*response.SalaryHistoryResponse, error) {
	// 1. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Criar array de 12 meses
	monthlySummaries := make([]response.MonthlySummary, 12)

	// 3. Para cada mês, buscar summary
	for month := 1; month <= 12; month++ {
		summary, err := uc.salaryEntryGateway.GetMonthlySummary(ctx, incomeID, month, year)
		if err != nil {
			return nil, fmt.Errorf("failed to get monthly summary for month %d: %w", month, err)
		}

		monthlySummaries[month-1] = response.MonthlySummary{
			Month:          month,
			BaseAmount:     summary.BaseAmount,
			TotalAdditions: summary.TotalAdditions,
			TotalDiscounts: summary.TotalDiscounts,
			NetAmount:      summary.NetAmount,
		}
	}

	// 4. Montar response
	history := &response.SalaryHistoryResponse{
		IncomeID:         incomeID,
		Description:      income.Description(),
		Year:             year,
		MonthlySummaries: monthlySummaries,
	}

	return history, nil
}

// RecalculateCLT recalculates CLT entries for a specific month
func (uc *useCaseImpl) RecalculateCLT(ctx context.Context, incomeID string, month, year int) error {
	// 1. Buscar income
	income, err := uc.incomeGateway.Get(ctx, incomeID)
	if err != nil {
		return fmt.Errorf("failed to get income: %w", err)
	}

	// 2. Validar que é CLT
	if income.DiscountMode() == nil || *income.DiscountMode() != "CLT" {
		return errors.New("income is not configured for CLT discounts")
	}

	// 3. Buscar todas as entries do período
	entries, err := uc.salaryEntryGateway.FindByIncomeAndPeriod(ctx, incomeID, month, year)
	if err != nil {
		return fmt.Errorf("failed to get entries: %w", err)
	}

	// 4. Calcular gross amount (base + additions)
	grossAmount := income.Amount()
	for _, entry := range entries {
		if entry.IsAddition() {
			grossAmount += entry.Amount()
		}
	}

	// 5. Deletar entries antigas de INSS e IRRF
	for _, entry := range entries {
		if entry.Category() == "INSS" || entry.Category() == "IRRF" {
			if err := uc.salaryEntryGateway.Delete(ctx, entry.ID()); err != nil {
				return fmt.Errorf("failed to delete old CLT entry: %w", err)
			}
		}
	}

	// 6. Calcular novos descontos CLT
	cltEntries, err := uc.calculator.CalculateCLTDiscounts(grossAmount, month, year)
	if err != nil {
		return fmt.Errorf("failed to calculate CLT discounts: %w", err)
	}

	// 7. Criar novas entries
	for _, cltEntry := range cltEntries {
		entryID := uuid.New().String()
		entry, err := models.NewSalaryEntry(
			entryID,
			incomeID,
			month,
			year,
			cltEntry.EntryType,
			cltEntry.Category,
			cltEntry.CalculationType,
			cltEntry.BaseValue,
			cltEntry.Multiplier,
			cltEntry.Amount,
			&cltEntry.Description,
			cltEntry.ReferenceDate,
		)
		if err != nil {
			return fmt.Errorf("failed to create CLT entry model: %w", err)
		}

		if err := uc.salaryEntryGateway.Create(ctx, entry); err != nil {
			return fmt.Errorf("failed to save CLT entry: %w", err)
		}
	}

	return nil
}
