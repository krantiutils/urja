.PHONY: build run test lint clean docker-up docker-down tidy

BINARY_NAME := urja-api
BUILD_DIR := ./bin

# Build the API binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

# Run the API server locally
run:
	go run ./cmd/api

# Run all tests
test:
	go test -v -race -count=1 ./...

# Run tests with coverage
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Run go vet
vet:
	go vet ./...

# Tidy modules
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

# Start Docker services (PostgreSQL + Redis)
docker-up:
	docker compose up -d postgres redis

# Start all Docker services including the API
docker-up-all:
	docker compose up -d --build

# Stop Docker services
docker-down:
	docker compose down

# Stop Docker services and remove volumes
docker-down-clean:
	docker compose down -v
