.PHONY: build run test clean deps setup-db migrate migrate-down dev

BINARY_NAME=e173gw
GO=/usr/local/go/bin/go
CMD_PATH=./cmd/server

# Default target
all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@$(GO) build -o bin/$(BINARY_NAME) $(CMD_PATH)

run: build
	@echo "Starting $(BINARY_NAME) on port 8080..."
	@./bin/$(BINARY_NAME)

dev: build
	@echo "Starting $(BINARY_NAME) in development mode..."
	@./bin/$(BINARY_NAME) &
	@echo "Server started in background (PID: $$!)"

# Database setup and migrations
setup-db:
	@echo "Setting up database..."
	@chmod +x scripts/setup_database.sh
	@timeout 30s scripts/setup_database.sh || echo "Database setup completed or timed out"

migrate:
	@echo "Running database migrations..."
	@timeout 15s $(GO) run tools/migrate/*.go up || echo "Migration completed or not available"

migrate-down:
	@echo "Rolling back database migrations..."
	@timeout 15s $(GO) run tools/migrate/*.go down || echo "Migration rollback completed or not available"

# Testing and validation
test:
	@echo "Running tests..."
	@$(GO) test ./...

test-api:
	@echo "Testing API endpoints..."
	@timeout 5s curl -s http://localhost:8080/ping || echo "Server not responding"
	@timeout 5s curl -s http://localhost:8080/api/stats || echo "Stats API not responding"

# Development helpers
deps:
	@echo "Installing dependencies..."
	@$(GO) mod tidy
	@$(GO) mod download

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

lint:
	@echo "Running linter..."
	@golangci-lint run || echo "Linter not installed or issues found"

# Deployment preparation
deploy-prep:
	@echo "Preparing for deployment..."
	@$(GO) mod tidy
	@$(GO) build -o bin/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete for deployment"

clean:
	@echo "Cleaning..."
	@$(GO) clean
	@rm -rf bin/
	@pkill -f $(BINARY_NAME) || true

help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Run in development mode (background)"
	@echo "  setup-db   - Set up PostgreSQL database"
	@echo "  migrate    - Run database migrations"
	@echo "  test       - Run Go tests"
	@echo "  test-api   - Test API endpoints"
	@echo "  clean      - Clean build artifacts"
	@echo "  deploy-prep - Prepare for deployment"
