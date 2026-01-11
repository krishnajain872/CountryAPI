#!/usr/bin/env bash
# scripts/dev.sh - All-in-one Go Developer Script (cross-platform)
# Run from project root using: ./scripts/dev.sh <command>

set -e
set -o pipefail

# ========================================
# COLORS
# ========================================
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"
RESET="\033[0m"

# ========================================
# HELP
# ========================================
function help() {
    echo -e "${CYAN}Usage: ./scripts/dev.sh <command>${RESET}"
    echo -e "Commands:"
    echo -e "  help             Show this help"
    echo -e "  build            Build the project"
    echo -e "  run              Run the project"
    echo -e "  test             Run all tests"
    echo -e "  coverage         Run tests and generate coverage report"
    echo -e "  race             Run tests with race detector"
    echo -e "  fmt              Format the code"
    echo -e "  lint             Run linter"
    echo -e "  clean            Clean build and coverage artifacts"
    echo -e "  deps             Download dependencies"
    echo -e "  docker-build     Build Docker image"
    echo -e "  docker-run       Run Docker container"
}

# ========================================
# CHECK GO INSTALLED
# ========================================
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go to continue."
    exit 1
fi

# ========================================
# SET VARIABLES (relative to scripts/)
# ========================================
TARGET=$1
PROJECT_ROOT=$(dirname "$0")/..
BIN_DIR="${PROJECT_ROOT}/bin"
APP_NAME="server"
COVERAGE_FILE="${PROJECT_ROOT}/coverage.out"
COVERAGE_HTML="${PROJECT_ROOT}/coverage.html"

# ========================================
# COMMANDS
# ========================================
case $TARGET in
    build)
        echo -e "${YELLOW}🔨 Building project...${RESET}"
        mkdir -p ${BIN_DIR}
        go build -o ${BIN_DIR}/${APP_NAME} ${PROJECT_ROOT}/cmd/server/main.go
        echo -e "${GREEN}✅ Build complete: ${BIN_DIR}/${APP_NAME}${RESET}"
        ;;

    run)
        echo -e "${YELLOW}🚀 Running server...${RESET}"
        go run ${PROJECT_ROOT}/cmd/server/.
        ;;

    test)
        echo -e "${YELLOW}🧪 Running tests...${RESET}"
        go test -v ${PROJECT_ROOT}/...
        ;;

    coverage)
        echo -e "${YELLOW}📊 Running coverage...${RESET}"
        go test ${PROJECT_ROOT}/... -covermode=atomic -coverpkg=${PROJECT_ROOT}/... -coverprofile=${COVERAGE_FILE}
        go tool cover -func=${COVERAGE_FILE}
        go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}
        echo -e "${GREEN}✅ Coverage report: ${COVERAGE_HTML}${RESET}"
        ;;

    race)
        echo -e "${YELLOW}🏁 Running race detector...${RESET}"
        go test -race ${PROJECT_ROOT}/...
        ;;

    fmt)
        echo -e "${YELLOW}🧹 Formatting code...${RESET}"
        go fmt ${PROJECT_ROOT}/...
        goimports -w ${PROJECT_ROOT}
        ;;

    lint)
        echo -e "${YELLOW}🔍 Running linter...${RESET}"
        golangci-lint run ${PROJECT_ROOT}/...
        ;;

    clean)
        echo -e "${YELLOW}🧽 Cleaning build and coverage...${RESET}"
        rm -rf ${BIN_DIR} ${COVERAGE_FILE} ${COVERAGE_HTML}
        echo -e "${GREEN}✅ Clean complete${RESET}"
        ;;

    deps)
        echo -e "${YELLOW}📦 Installing dependencies...${RESET}"
        cd ${PROJECT_ROOT}
        go mod download
        go mod tidy
        ;;

    docker-build)
        echo -e "${YELLOW}🐳 Building Docker image...${RESET}"
        docker build -t country-search-api:latest ${PROJECT_ROOT}
        ;;

    docker-run)
        echo -e "${YELLOW}🐳 Running Docker container...${RESET}"
        docker run -p 8000:8000 country-search-api:latest
        ;;

    help|""|*)
        help
        ;;
esac
