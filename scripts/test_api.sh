#!/bin/bash
# ========================================
# FILE: scripts/run_all_tests.sh
# Test runner script (adjusted for current folder structure)
# ========================================

set -e

echo "=========================================="
echo "Running Complete Test Suite"
echo "=========================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
FAILED_TESTS=0

# -----------------------------
# Phase 1: Unit Tests
# -----------------------------
echo ""
echo "${YELLOW}Phase 1: Unit Tests${NC}"
echo "=========================================="
if go test -v ./tests/unit/...; then
    echo "${GREEN}Unit tests passed${NC}"
else
    echo "${RED}Unit tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 2: Integration Tests
# -----------------------------
echo ""
echo "${YELLOW}Phase 2: Integration Tests${NC}"
echo "=========================================="
if go test -v ./tests/integration/...; then
    echo "${GREEN}Integration tests passed${NC}"
else
    echo "${RED}Integration tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 3: Race Condition Tests (optional)
# -----------------------------
echo ""
echo "${YELLOW}Phase 3: Race Condition Tests${NC}"
echo "=========================================="
# Enable CGO for -race
export CGO_ENABLED=1
if go test -race -v ./tests/unit/...; then
    echo "${GREEN}Race tests passed${NC}"
else
    echo "${RED}Race tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 4: Coverage Report
# -----------------------------
echo ""
echo "${YELLOW}Phase 4: Coverage Report${NC}"
echo "=========================================="
go test ./tests/... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1

# -----------------------------
# Summary
# -----------------------------
echo ""
echo "=========================================="
if [ $FAILED_TESTS -eq 0 ]; then
    echo "${GREEN}✅ All test phases passed!${NC}"
else
    echo "${RED}❌ $FAILED_TESTS test phase(s) failed${NC}"
    exit 1
fi
echo "=========================================="
