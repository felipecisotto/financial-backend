-- Migration: Add salary tracking fields to incomes table
-- Description: Adds discount_mode and hourly_rate columns to enable salary entry tracking

-- Add new columns
ALTER TABLE incomes
ADD COLUMN discount_mode VARCHAR,
ADD COLUMN hourly_rate DECIMAL(10,2);

-- Add constraints
ALTER TABLE incomes
ADD CONSTRAINT check_discount_mode
CHECK (discount_mode IS NULL OR discount_mode IN ('CLT', 'custom'));

ALTER TABLE incomes
ADD CONSTRAINT check_hourly_rate_positive
CHECK (hourly_rate IS NULL OR hourly_rate > 0);

-- Rollback:
-- ALTER TABLE incomes DROP CONSTRAINT check_hourly_rate_positive;
-- ALTER TABLE incomes DROP CONSTRAINT check_discount_mode;
-- ALTER TABLE incomes DROP COLUMN hourly_rate;
-- ALTER TABLE incomes DROP COLUMN discount_mode;
