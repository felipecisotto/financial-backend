# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy .env file if it exists
RUN if [ -f /app/.env ]; then cp /app/.env .; fi

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
