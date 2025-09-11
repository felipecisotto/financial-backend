# Test commands
.PHONY: test test-unit test-integration test-coverage test-race install-test-deps clean-test fmt vet lint check

# Run all tests
test:
	go test ./... -v

# Run only unit tests (exclude integration)
test-unit:
	go test ./internal/controllers ./internal/usecases/... -v -short

# Run integration tests
test-integration:
	go test ./tests/integration/... -v

# Run tests with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	go test ./... -v -race

# Generate mocks
generate-mocks:
	go generate ./...

# Install test dependencies
install-test-deps:
	go install github.com/golang/mock/mockgen@latest
	go get -t ./...

# Clean test artifacts
clean-test:
	rm -f coverage.out coverage.html

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks (tests, lint, vet, fmt)
check: fmt vet test-coverage

# Development commands
.PHONY: run build docker-build docker-run deps tidy

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o main cmd/api/main.go

# Build Docker image
docker-build:
	docker build -t financial-backend .

# Run with Docker
docker-run:
	docker run -p 8080:8080 financial-backend

# Download dependencies
deps:
	go mod download

# Tidy dependencies
tidy:
	go mod tidy