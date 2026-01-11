.PHONY: help build run test test-coverage test-race clean lint fmt deps docker-build docker-run

APP_NAME := server
BIN_DIR := bin
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# ========================================
# HELP
# ========================================
help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ========================================
# BUILD & RUN
# ========================================
build: ## Build application
	@echo "🔨 Building..."
	@go build -o $(BIN_DIR)/$(APP_NAME) cmd/server/main.go
	@echo "✅ Build complete: $(BIN_DIR)/$(APP_NAME)"

run: ## Run application
	@echo "🚀 Running server..."
	@go run cmd/server/main.go

# ========================================
# TESTING
# ========================================
test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage: ## Run full coverage (cmd + internal + pkg)
	@echo "📊 Running coverage tests..."
	@go test ./... \
		-covermode=atomic \
		-coverpkg=./... \
		-coverprofile=$(COVERAGE_FILE)
	@go tool cover -func=$(COVERAGE_FILE)
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "✅ Coverage report generated: $(COVERAGE_HTML)"

test-race: ## Run race detector
	@echo "🏁 Running race detector..."
	@go test -race ./...

# ========================================
# QUALITY
# ========================================
fmt: ## Format code
	@echo "🧹 Formatting..."
	@go fmt ./...
	@goimports -w .

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

deps: ## Download & tidy dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

# ========================================
# CLEAN
# ========================================
clean: ## Clean artifacts
	@echo "🧽 Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "✅ Clean complete"

# ========================================
# DOCKER
# ========================================
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t country-search-api:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8000:8000 country-search-api:latest
