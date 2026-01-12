#!/usr/bin/env bash
# ========================================
# FILE: scripts/build.sh
# Complete Development Script with Fixed Coverage
# ========================================

set -e
set -o pipefail

# ========================================
# COLORS
# ========================================
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"
RED="\033[0;31m"
RESET="\033[0m"

# ========================================
# VARIABLES
# ========================================
TARGET=$1
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${PROJECT_ROOT}/bin"
APP_NAME="server"
COVERAGE_FILE="${PROJECT_ROOT}/coverage.out"
COVERAGE_HTML="${PROJECT_ROOT}/coverage.html"
COVERAGE_SUMMARY="${PROJECT_ROOT}/coverage_summary.txt"

# ========================================
# HELP
# ========================================
function help() {
    echo -e "${CYAN}========================================${RESET}"
    echo -e "${CYAN}Country Search API - Build Script${RESET}"
    echo -e "${CYAN}========================================${RESET}"
    echo -e ""
    echo -e "${YELLOW}Usage:${RESET} ./scripts/build.sh <command>"
    echo -e ""
    echo -e "${GREEN}Build & Run:${RESET}"
    echo -e "  ${CYAN}build${RESET}            Build the project"
    echo -e "  ${CYAN}run${RESET}              Run the server"
    echo -e "  ${CYAN}clean${RESET}            Clean build artifacts and coverage files"
    echo -e ""
    echo -e "${GREEN}Testing:${RESET}"
    echo -e "  ${CYAN}test${RESET}             Run all tests (unit + integration)"
    echo -e "  ${CYAN}test-unit${RESET}        Run only unit tests"
    echo -e "  ${CYAN}test-integration${RESET} Run only integration tests"
    echo -e "  ${CYAN}test-all${RESET}         Run comprehensive test suite with phases"
    echo -e "  ${CYAN}race${RESET}             Run tests with race detector"
    echo -e "  ${CYAN}race-all${RESET}         Run race detector on all test files"
    echo -e ""
    echo -e "${GREEN}Coverage:${RESET}"
    echo -e "  ${CYAN}coverage${RESET}         Generate full coverage report (HTML + Summary)"
    echo -e "  ${CYAN}coverage-unit${RESET}    Coverage for unit tests only"
    echo -e "  ${CYAN}coverage-pkg${RESET}     Coverage by package breakdown"
    echo -e ""
    echo -e "${GREEN}Code Quality:${RESET}"
    echo -e "  ${CYAN}fmt${RESET}              Format code (go fmt + goimports)"
    echo -e "  ${CYAN}lint${RESET}             Run linter (golangci-lint)"
    echo -e "  ${CYAN}vet${RESET}              Run go vet"
    echo -e ""
    echo -e "${GREEN}Dependencies:${RESET}"
    echo -e "  ${CYAN}deps${RESET}             Download and tidy dependencies"
    echo -e "  ${CYAN}deps-update${RESET}      Update all dependencies"
    echo -e ""
    echo -e "${GREEN}Docker:${RESET}"
    echo -e "  ${CYAN}docker-build${RESET}     Build Docker image"
    echo -e "  ${CYAN}docker-run${RESET}       Run Docker container"
    echo -e ""
    echo -e "${GREEN}Utilities:${RESET}"
    echo -e "  ${CYAN}help${RESET}             Show this help message"
    echo -e "  ${CYAN}version${RESET}          Show Go version"
    echo -e ""
    echo -e "${CYAN}========================================${RESET}"
}

# ========================================
# CHECK GO INSTALLED
# ========================================
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go to continue.${RESET}"
    exit 1
fi

# ========================================
# UTILITY FUNCTIONS
# ========================================
function print_header() {
    echo -e ""
    echo -e "${CYAN}==========================================${RESET}"
    echo -e "${CYAN}$1${RESET}"
    echo -e "${CYAN}==========================================${RESET}"
}

function print_success() {
    echo -e "${GREEN}✅ $1${RESET}"
}

function print_error() {
    echo -e "${RED}❌ $1${RESET}"
}

function print_info() {
    echo -e "${YELLOW}ℹ️  $1${RESET}"
}

# ========================================
# MAIN COMMANDS
# ========================================

case $TARGET in
    # ========================================
    # BUILD & RUN
    # ========================================
    build)
        print_header "Building Project"
        mkdir -p "${BIN_DIR}"
        cd "${PROJECT_ROOT}"
        go build -o "${BIN_DIR}/${APP_NAME}" ./cmd/server/
        print_success "Build complete: ${BIN_DIR}/${APP_NAME}"
        ;;

    run)
        print_header "Running Server"
        cd "${PROJECT_ROOT}"
        go run ./cmd/server/
        ;;

    clean)
        print_header "Cleaning Artifacts"
        rm -rf "${BIN_DIR}" "${COVERAGE_FILE}" "${COVERAGE_HTML}" "${COVERAGE_SUMMARY}"
        rm -f "${PROJECT_ROOT}"/*.out "${PROJECT_ROOT}"/*.html
        rm -rf "${PROJECT_ROOT}/reports"
        print_success "Clean complete"
        ;;

    # ========================================
    # TESTING COMMANDS
    # ========================================
    test)
        print_header "Running All Tests"
        cd "${PROJECT_ROOT}"
        go test -v ./tests/...
        print_success "All tests completed"
        ;;

    test-unit)
        print_header "Running Unit Tests"
        cd "${PROJECT_ROOT}"
        go test -v ./tests/unit/...
        print_success "Unit tests completed"
        ;;

    test-integration)
        print_header "Running Integration Tests"
        cd "${PROJECT_ROOT}"
        print_info "Integration tests require internet connection"
        go test -v ./tests/integration/... -timeout 30s || print_error "Integration tests failed (may need internet)"
        ;;

    test-all)
        print_header "Running Comprehensive Test Suite"
        cd "${PROJECT_ROOT}"
        
        TOTAL_PHASES=0
        FAILED_PHASES=0
        
        # Phase 1: Unit Tests
        echo ""
        echo -e "${YELLOW}📦 Phase 1: Unit Tests${RESET}"
        echo "=========================================="
        if go test -v ./tests/unit/...; then
            print_success "Unit tests passed"
        else
            print_error "Unit tests failed"
            ((FAILED_PHASES++))
        fi
        ((TOTAL_PHASES++))
        
        # Phase 2: Cache Layer Tests
        echo ""
        echo -e "${YELLOW}💾 Phase 2: Cache Layer Tests${RESET}"
        echo "=========================================="
        if go test -v ./pkg/cache/...; then
            print_success "Cache tests passed"
        else
            print_error "Cache tests failed"
            ((FAILED_PHASES++))
        fi
        ((TOTAL_PHASES++))
        
        # Phase 3: Service Layer Tests
        echo ""
        echo -e "${YELLOW}⚙️  Phase 3: Service Layer Tests${RESET}"
        echo "=========================================="
        if go test -v ./internal/service/...; then
            print_success "Service tests passed"
        else
            print_error "Service tests failed"
            ((FAILED_PHASES++))
        fi
        ((TOTAL_PHASES++))
        
        # Phase 4: Handler Layer Tests
        echo ""
        echo -e "${YELLOW}🌐 Phase 4: Handler Layer Tests${RESET}"
        echo "=========================================="
        if go test -v ./internal/handler/...; then
            print_success "Handler tests passed"
        else
            print_error "Handler tests failed"
            ((FAILED_PHASES++))
        fi
        ((TOTAL_PHASES++))
        
        # Phase 5: Integration Tests
        echo ""
        echo -e "${YELLOW}🔗 Phase 5: Integration Tests${RESET}"
        echo "=========================================="
        if go test -v ./tests/integration/... -timeout 30s; then
            print_success "Integration tests passed"
        else
            print_error "Integration tests failed (may need internet)"
            ((FAILED_PHASES++))
        fi
        ((TOTAL_PHASES++))
        
        # Summary
        echo ""
        print_header "Test Summary"
        echo -e "Total Phases: ${TOTAL_PHASES}"
        echo -e "${GREEN}Passed: $((TOTAL_PHASES - FAILED_PHASES))${RESET}"
        echo -e "${RED}Failed: ${FAILED_PHASES}${RESET}"
        
        if [ $FAILED_PHASES -eq 0 ]; then
            echo ""
            print_success "All test phases passed! 🎉"
            exit 0
        else
            echo ""
            print_error "$FAILED_PHASES test phase(s) failed"
            exit 1
        fi
        ;;

    # ========================================
    # RACE DETECTION
    # ========================================
    race)
        print_header "Running Race Detector (Unit Tests)"
        cd "${PROJECT_ROOT}"
        export CGO_ENABLED=1
        go test -race -v ./tests/unit/... -run "Race"
        print_success "Race detector tests completed"
        ;;

    race-all)
        print_header "Running Race Detector (All Tests)"
        cd "${PROJECT_ROOT}"
        export CGO_ENABLED=1
        
        echo ""
        echo -e "${YELLOW}🏁 Testing: Unit Tests${RESET}"
        go test -race -v ./tests/unit/... -run "Race" || print_error "Unit tests have race conditions"
        
        echo ""
        echo -e "${YELLOW}🏁 Testing: Cache Layer${RESET}"
        go test -race -v ./pkg/cache/... || print_error "Cache has race conditions"
        
        echo ""
        echo -e "${YELLOW}🏁 Testing: Service Layer${RESET}"
        go test -race -v ./internal/service/... || print_error "Service has race conditions"
        
        echo ""
        echo -e "${YELLOW}🏁 Testing: Handler Layer${RESET}"
        go test -race -v ./internal/handler/... || print_error "Handler has race conditions"
        
        print_success "Race detector completed on all layers"
        ;;

    # ========================================
    # COVERAGE COMMANDS (FIXED)
    # ========================================
    coverage)
        print_header "Generating Coverage Report"
        cd "${PROJECT_ROOT}"
        
        # Run tests with coverage - THE KEY FIX
        print_info "Running tests with coverage tracking..."
        
        # Use -coverpkg to specify which packages to track coverage for
        # This is CRITICAL when tests are in a separate directory
        go test ./tests/... \
            -coverprofile="${COVERAGE_FILE}" \
            -covermode=atomic \
            -coverpkg=./internal/...,./pkg/...,./cmd/... \
            -timeout=30s
        
        if [ ! -f "${COVERAGE_FILE}" ]; then
            print_error "Coverage file not generated"
            exit 1
        fi
        
        # Generate summary
        print_info "Generating coverage summary..."
        go tool cover -func="${COVERAGE_FILE}" > "${COVERAGE_SUMMARY}"
        
        # Print summary to console
        echo ""
        echo -e "${YELLOW}Coverage Summary by Package:${RESET}"
        echo "=========================================="
        go tool cover -func="${COVERAGE_FILE}" | grep -E '^github.com' | head -30
        
        echo ""
        echo -e "${YELLOW}Total Coverage:${RESET}"
        echo "=========================================="
        TOTAL=$(go tool cover -func="${COVERAGE_FILE}" | grep "total:" | awk '{print $3}')
        echo -e "${GREEN}${TOTAL}${RESET}"
        
        # Generate HTML
        print_info "Generating HTML coverage report..."
        go tool cover -html="${COVERAGE_FILE}" -o "${COVERAGE_HTML}"
        
        if [ -f "${COVERAGE_HTML}" ]; then
            echo ""
            print_success "Coverage reports generated:"
            echo -e "  📄 Summary: ${COVERAGE_SUMMARY}"
            echo -e "  📊 HTML:    ${COVERAGE_HTML}"
            echo ""
            print_info "Opening ${COVERAGE_HTML} in browser..."
            
            # Open in browser
            if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
                start "${COVERAGE_HTML}" 2>/dev/null || true
            elif [[ "$OSTYPE" == "darwin"* ]]; then
                open "${COVERAGE_HTML}" 2>/dev/null || true
            else
                xdg-open "${COVERAGE_HTML}" 2>/dev/null || true
            fi
        else
            print_error "Failed to generate HTML report"
            exit 1
        fi
        ;;

    coverage-unit)
        print_header "Generating Unit Test Coverage"
        cd "${PROJECT_ROOT}"
        
        go test ./tests/unit/... \
            -coverprofile="${COVERAGE_FILE}" \
            -covermode=atomic \
            -coverpkg=./internal/...,./pkg/...,./cmd/...
        
        go tool cover -func="${COVERAGE_FILE}" | tail -n 1
        go tool cover -html="${COVERAGE_FILE}" -o "${COVERAGE_HTML}"
        
        print_success "Unit test coverage: ${COVERAGE_HTML}"
        
        # Open in browser
        if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
            start "${COVERAGE_HTML}" 2>/dev/null || true
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            open "${COVERAGE_HTML}" 2>/dev/null || true
        else
            xdg-open "${COVERAGE_HTML}" 2>/dev/null || true
        fi
        ;;

    coverage-pkg)
        print_header "Coverage by Package"
        cd "${PROJECT_ROOT}"
        
        go test ./tests/... \
            -coverprofile="${COVERAGE_FILE}" \
            -covermode=atomic \
            -coverpkg=./internal/...,./pkg/...,./cmd/...
        
        echo ""
        echo -e "${YELLOW}Coverage Breakdown:${RESET}"
        echo "=========================================="
        
        # Cache layer
        echo -e "\n${CYAN}📦 Cache Layer:${RESET}"
        go tool cover -func="${COVERAGE_FILE}" | grep "pkg/cache" || echo "No coverage data"
        
        # Service layer
        echo -e "\n${CYAN}⚙️  Service Layer:${RESET}"
        go tool cover -func="${COVERAGE_FILE}" | grep "internal/service" || echo "No coverage data"
        
        # Handler layer
        echo -e "\n${CYAN}🌐 Handler Layer:${RESET}"
        go tool cover -func="${COVERAGE_FILE}" | grep "internal/handler" || echo "No coverage data"
        
        # Repository layer
        echo -e "\n${CYAN}💾 Repository Layer:${RESET}"
        go tool cover -func="${COVERAGE_FILE}" | grep "internal/repository" || echo "No coverage data"
        
        # Domain layer
        echo -e "\n${CYAN}📋 Domain Layer:${RESET}"
        go tool cover -func="${COVERAGE_FILE}" | grep "internal/domain" || echo "No coverage data"
        
        # Total
        echo ""
        echo -e "${YELLOW}Total Coverage:${RESET}"
        echo "=========================================="
        go tool cover -func="${COVERAGE_FILE}" | grep "total:"
        ;;

    # ========================================
    # CODE QUALITY
    # ========================================
    fmt)
        print_header "Formatting Code"
        cd "${PROJECT_ROOT}"
        
        print_info "Running go fmt..."
        go fmt ./...
        
        if command -v goimports &> /dev/null; then
            print_info "Running goimports..."
            goimports -w .
        else
            print_info "goimports not installed, skipping..."
            print_info "Install with: go install golang.org/x/tools/cmd/goimports@latest"
        fi
        
        print_success "Code formatted"
        ;;

    lint)
        print_header "Running Linter"
        cd "${PROJECT_ROOT}"
        
        if command -v golangci-lint &> /dev/null; then
            golangci-lint run ./...
            print_success "Linting complete"
        else
            print_error "golangci-lint not installed"
            print_info "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
            exit 1
        fi
        ;;

    vet)
        print_header "Running Go Vet"
        cd "${PROJECT_ROOT}"
        go vet ./...
        print_success "Vet checks passed"
        ;;

    # ========================================
    # DEPENDENCIES
    # ========================================
    deps)
        print_header "Managing Dependencies"
        cd "${PROJECT_ROOT}"
        
        print_info "Downloading dependencies..."
        go mod download
        
        print_info "Tidying dependencies..."
        go mod tidy
        
        print_success "Dependencies updated"
        ;;

    deps-update)
        print_header "Updating All Dependencies"
        cd "${PROJECT_ROOT}"
        
        print_info "Updating dependencies..."
        go get -u ./...
        go mod tidy
        
        print_success "Dependencies updated to latest versions"
        ;;

    # ========================================
    # DOCKER
    # ========================================
    docker-build)
        print_header "Building Docker Image"
        cd "${PROJECT_ROOT}"
        docker build -t country-search-api:latest .
        print_success "Docker image built: country-search-api:latest"
        ;;

    docker-run)
        print_header "Running Docker Container"
        docker run -p 8000:8000 country-search-api:latest
        ;;

    # ========================================
    # UTILITIES
    # ========================================
    version)
        print_header "Go Version Information"
        go version
        echo ""
        echo "Project Root: ${PROJECT_ROOT}"
        echo "Go Modules: $(go env GOMOD)"
        ;;

    # ========================================
    # DEFAULT / HELP
    # ========================================
    help|""|*)
        help
        ;;
esac