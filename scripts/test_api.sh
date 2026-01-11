
// ========================================
// FILE: scripts/run_all_tests.sh
// Test runner script
// ========================================

#!/bin/bash

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
PASSED_TESTS=0
FAILED_TESTS=0

echo ""
echo "${YELLOW}Phase 1: Cache Layer Tests${NC}"
echo "=========================================="
go test -v ./pkg/cache/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 2: Domain Layer Tests${NC}"
echo "=========================================="
go test -v ./internal/domain/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 3: Service Layer Tests${NC}"
echo "=========================================="
go test -v ./internal/service/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 4: Handler Layer Tests${NC}"
echo "=========================================="
go test -v ./internal/handler/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 5: Repository Layer Tests${NC}"
echo "=========================================="
go test -v ./internal/repository/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 6: Race Condition Tests${NC}"
echo "=========================================="
go test -race -v ./pkg/cache/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 7: Integration Tests${NC}"
echo "=========================================="
go test -v ./tests/integration/... || FAILED_TESTS=$((FAILED_TESTS+1))
TOTAL_TESTS=$((TOTAL_TESTS+1))

echo ""
echo "${YELLOW}Phase 8: Coverage Report${NC}"
echo "=========================================="
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1

echo ""
echo "=========================================="
if [ $FAILED_TESTS -eq 0 ]; then
    echo "${GREEN}✅ All test phases passed!${NC}"
else
    echo "${RED}❌ $FAILED_TESTS test phase(s) failed${NC}"
    exit 1
fi

echo "=========================================="
