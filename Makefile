.PHONY: help build run test test-coverage test-race clean lint fmt

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

run: ## Run the application
	@echo "Running server..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test -race ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t country-search-api:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8000:8000 country-search-api:latest

# ========================================
# FILE: .env.example
# ========================================
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8000
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# Cache Configuration
CACHE_TTL=5m
CACHE_MAX_SIZE=1000

# External API Configuration
REST_COUNTRIES_API_URL=https://restcountries.com/v3.1
API_TIMEOUT=10s

# Logger Configuration
LOG_LEVEL=info
LOG_FORMAT=json
