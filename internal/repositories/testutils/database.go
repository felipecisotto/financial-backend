package testutils

import (
	"context"
	"testing"

	"financial-backend/internal/domain/entities"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDatabase creates a PostgreSQL test container and returns a connected GORM DB instance
func SetupTestDatabase(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Start PostgreSQL container
	postgresC, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("financial_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get connection string
	connStr, err := postgresC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect with GORM
	db, err := gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate all entities
	err = db.AutoMigrate(
		&entities.Expense{},
		&entities.Budget{},
		&entities.Income{},
		&entities.BudgetMovement{},
	)
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	cleanup := func() {
		if err := postgresC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate PostgreSQL container: %v", err)
		}
	}

	return db, cleanup
}
