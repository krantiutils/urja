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

# ─── Database Migrations (golang-migrate) ───────────────────

DB_URL ?= postgres://urja:urja@localhost:5432/urja?sslmode=disable
MIGRATIONS_DIR = db/migrations

MIGRATE = migrate

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create migrate-reset

## Run all pending migrations
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## Rollback the last migration
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## Rollback all migrations
migrate-down-all:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

## Force a specific migration version (use when dirty)
migrate-force:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

## Show current migration version
migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

## Create a new migration (usage: make migrate-create NAME=create_foo)
migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

## Reset database (down all, then up all)
migrate-reset:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
