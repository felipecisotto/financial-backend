# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a financial backend API built with Go, using Clean Architecture principles. It manages personal financial data including expenses, income, budgets, and budget movements. The application uses Gin web framework, GORM for database operations, PostgreSQL as the database, and includes OpenTelemetry for observability.

## Common Commands

### Development
```bash
# Run the application
go run cmd/api/main.go

# Build the application
go build -o main cmd/api/main.go

# Run with Docker
docker build -t financial-backend .
docker run -p 8080:8080 financial-backend

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Format code
go fmt ./...

# Vet code
go vet ./...
```

### Testing
The project has comprehensive test coverage with unit tests, integration tests, and CI/CD pipeline:
```bash
# Run all tests
make test

# Run only unit tests (controllers and use cases)
make test-unit

# Run tests with coverage report
make test-coverage

# Run tests with race detection
make test-race

# Install testing dependencies
make install-test-deps

# Clean test artifacts
make clean-test

# Run all quality checks (format, vet, lint, test)
make check
```

## Architecture

The project follows Clean Architecture with clear separation of concerns:

### Directory Structure
- `cmd/api/` - Application entry point and main configuration
- `internal/` - Private application code
  - `controllers/` - HTTP handlers (Gin controllers)
  - `usecases/` - Business logic layer
  - `gateways/` - Interface adapters for repositories
  - `repositories/` - Data access layer implementations
  - `entities/` - Core business entities
  - `models/` - Data transfer objects and view models
  - `dtos/` - Data transfer objects
  - `views/` - Response view models
  - `mappers/` - Object mapping utilities
  - `events/` - Event handling system
- `pkg/` - Reusable packages
  - `config/` - Configuration management and database setup
  - `telemetry/` - OpenTelemetry instrumentation

### Key Components

**Main Modules:**
- **Expenses** - Expense tracking and management
- **Income** - Income recording and management  
- **Budget** - Budget creation and management
- **Budget Movement** - Budget allocation and movement tracking
- **Dashboard** - Financial overview and analytics

**Technology Stack:**
- **Web Framework:** Gin (github.com/gin-gonic/gin)
- **ORM:** GORM with PostgreSQL driver
- **Database:** PostgreSQL
- **Observability:** OpenTelemetry with Prometheus metrics
- **Configuration:** Environment variables with godotenv

### Architecture Flow
1. HTTP requests hit Controllers (Gin handlers)
2. Controllers call UseCases (business logic)
3. UseCases interact with Gateways (interface adapters)
4. Gateways use Repositories for data persistence
5. Event system handles cross-cutting concerns (e.g., budget movements on expense creation)

## Configuration

The application uses environment variables for configuration:
- `SERVER_ADDRESS` - Server listening address (default: ":8080")
- `DB_HOST` - Database host (default: "localhost")
- `DB_PORT` - Database port (default: "5432") 
- `DB_USER` - Database username (default: "postgres")
- `DB_PASSWORD` - Database password (default: "postgres")
- `DB_NAME` - Database name (default: "financial")
- `DEFAULT_DUE_DATE` - Default due date for expenses (default: "15")

Configuration is loaded from environment variables with `.env` file support.

## Database

- **Database:** PostgreSQL
- **ORM:** GORM with auto-migration enabled
- **Entities:** Budget, Expense, Income, BudgetMovement
- Database connection is managed as a singleton in `pkg/config/config.go`

## Development Guidelines

### Code Organization
- Follow the existing Clean Architecture pattern
- Keep business logic in UseCases
- Use Gateways as interfaces between UseCases and Repositories
- Controllers should be thin, only handling HTTP concerns
- Use the existing event system for cross-module communication

### Adding New Features
1. Define entities in `internal/entities/`
2. Create repository interface and implementation in `internal/repositories/{feature}/`
3. Create gateway in `internal/gateways/`
4. Implement use case in `internal/usecases/{feature}/`
5. Add controller in `internal/controllers/`
6. Register routes in the controller's `RegisterRoutes` method
7. Wire everything together in `cmd/api/main.go`

### Existing Patterns
- Dependency injection is done manually in main.go
- Error handling follows Go conventions
- All controllers have a `RegisterRoutes(router *gin.RouterGroup)` method
- Database models use GORM conventions
- Event-driven architecture for cross-module communication

## API Structure

The API follows REST conventions with the following structure:
- Base path: `/api`
- Health check: `/ping`
- CORS is enabled for all origins
- OpenTelemetry middleware is configured for observability

## Deployment

- **Docker:** Multi-stage build with Alpine Linux
- **CI/CD:** GitHub Actions workflow for automated Docker builds
- **Registry:** GitHub Container Registry (ghcr.io)
- **Versioning:** Automatic semantic versioning based on branch names (feature/* = minor, fix/* = patch)