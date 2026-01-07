-- Migration: Create salary_entries table
-- Description: Creates table to track salary entries (discounts and additions) for monthly income

CREATE TABLE salary_entries (
    id VARCHAR PRIMARY KEY,
    income_id VARCHAR NOT NULL REFERENCES incomes(id) ON DELETE CASCADE,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INT NOT NULL CHECK (year >= 2020),

    entry_type VARCHAR NOT NULL CHECK (entry_type IN ('discount', 'addition')),
    category VARCHAR NOT NULL,
    calculation_type VARCHAR NOT NULL CHECK (calculation_type IN ('fixed', 'percentage', 'computed')),

    base_value DECIMAL(10,2) NOT NULL DEFAULT 0,
    multiplier DECIMAL(10,4) NOT NULL DEFAULT 1.0,
    amount DECIMAL(10,2) NOT NULL,

    description TEXT,
    reference_date DATE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for optimized queries
CREATE INDEX idx_salary_entries_income_period ON salary_entries(income_id, month, year);
CREATE INDEX idx_salary_entries_type ON salary_entries(entry_type);
CREATE INDEX idx_salary_entries_category ON salary_entries(category);
CREATE INDEX idx_salary_entries_date ON salary_entries(reference_date);

-- Rollback:
-- DROP INDEX IF EXISTS idx_salary_entries_date;
-- DROP INDEX IF EXISTS idx_salary_entries_category;
-- DROP INDEX IF EXISTS idx_salary_entries_type;
-- DROP INDEX IF EXISTS idx_salary_entries_income_period;
-- DROP TABLE IF EXISTS salary_entries;
