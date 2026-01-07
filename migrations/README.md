# Database Migrations

This directory contains SQL migration files for the salary entries feature.

## Migration Files

1. **001_add_salary_fields_to_incomes.sql** - Adds salary tracking fields to incomes table
   - Adds `discount_mode` column (VARCHAR, nullable, values: 'CLT' or 'custom')
   - Adds `hourly_rate` column (DECIMAL(10,2), nullable, must be > 0)

2. **002_create_salary_entries.sql** - Creates salary_entries table
   - Tracks all salary additions and discounts for a given month/year
   - Supports multiple calculation types (fixed, percentage, computed)
   - Includes indexes for optimized queries

## How to Apply

This project uses GORM with AutoMigrate, which automatically creates and updates tables based on entity definitions. These SQL files are provided for:

1. **Documentation** - Understanding the database schema
2. **Manual Migration** - If needed for production environments
3. **Rollback Reference** - Each file includes rollback commands in comments

### Automatic Migration (Default)

The application automatically applies migrations on startup via GORM AutoMigrate in `pkg/config/config.go`.

### Manual Migration (Optional)

If you need to apply migrations manually:

```bash
# Connect to PostgreSQL
psql -h localhost -U postgres -d financial

# Apply migrations in order
\i migrations/001_add_salary_fields_to_incomes.sql
\i migrations/002_create_salary_entries.sql
```

## Rollback

To rollback migrations, use the commented rollback commands at the end of each migration file in reverse order:

```bash
# First, rollback 002
# Then, rollback 001
```

## Notes

- Migrations should be applied in numerical order
- All migrations include CHECK constraints for data validation
- Indexes are created for commonly queried columns
- Foreign key constraints ensure referential integrity
