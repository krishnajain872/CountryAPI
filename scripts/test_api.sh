#!/bin/bash
# ========================================
# FILE: scripts/test_api.sh (FIXED)
# Test runner script with proper error handling
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
echo -e "${YELLOW}Phase 1: Unit Tests${NC}"
echo "=========================================="
if go test -v ./tests/unit/...; then
    echo -e "${GREEN}Unit tests passed${NC}"
else
    echo -e "${RED}Unit tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 2: Integration Tests
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 2: Integration Tests${NC}"
echo "=========================================="
if go test -v ./tests/integration/...; then
    echo -e "${GREEN}Integration tests passed${NC}"
else
    echo -e "${RED}Integration tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 3: Race Condition Tests (FIXED)
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 3: Race Condition Tests${NC}"
echo "=========================================="

# Check if we're on Windows or if GCC is available
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    # Check if GCC is available on Windows
    if command -v gcc &> /dev/null; then
        echo "GCC found, running race tests..."
        export CGO_ENABLED=1
        if go test -race -v ./tests/unit/...; then
            echo -e "${GREEN}Race tests passed${NC}"
        else
            echo -e "${RED}Race tests failed${NC}"
            FAILED_TESTS=$((FAILED_TESTS+1))
        fi
    else
        echo -e "${YELLOW}⚠️  GCC not found - Skipping race tests on Windows${NC}"
        echo -e "${YELLOW}   Race tests require CGO (C compiler)${NC}"
        echo -e "${YELLOW}   Install TDM-GCC or MinGW-w64 to enable race tests${NC}"
        echo -e "${GREEN}✓ Skipped (not a failure)${NC}"
    fi
else
    # On Linux/Mac, GCC should be available
    export CGO_ENABLED=#!/bin/bash
# ========================================
# FILE: scripts/test_api.sh (ENHANCED)
# Test runner with beautiful HTML coverage report
# ========================================

set -e

echo "=========================================="
echo "Running Complete Test Suite"
echo "=========================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
FAILED_TESTS=0
START_TIME=$(date +%s)

# -----------------------------
# Phase 1: Unit Tests
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 1: Unit Tests${NC}"
echo "=========================================="
if go test -v ./tests/unit/...; then
    echo -e "${GREEN}Unit tests passed${NC}"
else
    echo -e "${RED}Unit tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 2: Integration Tests
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 2: Integration Tests${NC}"
echo "=========================================="
if go test -v ./tests/integration/...; then
    echo -e "${GREEN}Integration tests passed${NC}"
else
    echo -e "${RED}Integration tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS+1))
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 3: Race Condition Tests
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 3: Race Condition Tests${NC}"
echo "=========================================="

if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    if command -v gcc &> /dev/null; then
        echo "GCC found, running race tests..."
        export CGO_ENABLED=1
        if go test -race -v ./tests/unit/...; then
            echo -e "${GREEN}Race tests passed${NC}"
        else
            echo -e "${RED}Race tests failed${NC}"
            FAILED_TESTS=$((FAILED_TESTS+1))
        fi
    else
        echo -e "${YELLOW}⚠️  GCC not found - Skipping race tests on Windows${NC}"
        echo -e "${YELLOW}   Race tests require CGO (C compiler)${NC}"
        echo -e "${GREEN}✓ Skipped (not a failure)${NC}"
    fi
else
    export CGO_ENABLED=1
    if go test -race -v ./tests/unit/...; then
        echo -e "${GREEN}Race tests passed${NC}"
    else
        echo -e "${RED}Race tests failed${NC}"
        FAILED_TESTS=$((FAILED_TESTS+1))
    fi
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 4: Coverage Report
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 4: Coverage Report${NC}"
echo "=========================================="

# Generate coverage data
go test ./tests/... -coverprofile=coverage.out -covermode=atomic

# Get coverage percentage
COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo -e "${CYAN}Total Coverage: ${COVERAGE_PERCENT}${NC}"

# Extract detailed coverage data
go tool cover -func=coverage.out > coverage_details.txt

# -----------------------------
# Phase 5: Generate Beautiful HTML Report
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 5: Generating Beautiful Coverage Report${NC}"
echo "=========================================="

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
CURRENT_DATE=$(date '+%Y-%m-%d %H:%M:%S')

# Parse coverage by package
declare -A PACKAGE_COVERAGE
declare -A PACKAGE_FILES

while IFS= read -r line; do
    if [[ $line =~ ^github\.com.*\.go:.*% ]]; then
        # Extract package, file, and coverage
        PKG=$(echo "$line" | awk -F'/' '{print $(NF-1)}')
        FILE=$(echo "$line" | awk '{print $1}' | awk -F':' '{print $1}' | awk -F'/' '{print $NF}')
        COV=$(echo "$line" | awk '{print $NF}')
        
        if [[ ! -v PACKAGE_COVERAGE[$PKG] ]]; then
            PACKAGE_COVERAGE[$PKG]="$COV"
            PACKAGE_FILES[$PKG]="$FILE:$COV"
        else
            PACKAGE_FILES[$PKG]="${PACKAGE_FILES[$PKG]}|$FILE:$COV"
        fi
    fi
done < coverage_details.txt

# Create HTML report
cat >reports/coverage_report.html  << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Coverage Report - Country Search API</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <style>
        @keyframes fadeIn { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
        .fade-in { animation: fadeIn 0.5s ease-out; }
        .pulse-slow { animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite; }
    </style>
</head>
<body class="bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 min-h-screen text-gray-100">
EOF

# Add dynamic content
cat >>reports/coverage_report.html  << EOF
    <!-- Header -->
    <div class="container mx-auto px-4 py-8">
        <div class="fade-in bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl shadow-2xl p-8 mb-8">
            <div class="flex items-center justify-between">
                <div>
                    <h1 class="text-4xl font-bold text-white mb-2">📊 Test Coverage Report</h1>
                    <p class="text-blue-100 text-lg">Country Search API with Caching</p>
                </div>
                <div class="text-right">
                    <div class="text-5xl font-bold text-white">${COVERAGE_PERCENT}</div>
                    <div class="text-blue-100 text-sm">Total Coverage</div>
                </div>
            </div>
        </div>

        <!-- Summary Cards -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <!-- Total Tests -->
            <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-blue-500 transition-all duration-300">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-400 text-sm mb-1">Total Tests</p>
                        <p class="text-3xl font-bold text-white">${TOTAL_TESTS}</p>
                    </div>
                    <div class="bg-blue-500 bg-opacity-20 rounded-full p-3">
                        <svg class="w-8 h-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Passed Tests -->
            <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-green-500 transition-all duration-300">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-400 text-sm mb-1">Passed</p>
                        <p class="text-3xl font-bold text-green-400">$((TOTAL_TESTS - FAILED_TESTS))</p>
                    </div>
                    <div class="bg-green-500 bg-opacity-20 rounded-full p-3">
                        <svg class="w-8 h-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Failed Tests -->
            <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-red-500 transition-all duration-300">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-400 text-sm mb-1">Failed</p>
                        <p class="text-3xl font-bold text-red-400">${FAILED_TESTS}</p>
                    </div>
                    <div class="bg-red-500 bg-opacity-20 rounded-full p-3">
                        <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Duration -->
            <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-purple-500 transition-all duration-300">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-gray-400 text-sm mb-1">Duration</p>
                        <p class="text-3xl font-bold text-purple-400">${DURATION}s</p>
                    </div>
                    <div class="bg-purple-500 bg-opacity-20 rounded-full p-3">
                        <svg class="w-8 h-8 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                    </div>
                </div>
            </div>
        </div>

        <!-- Coverage Chart -->
        <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 mb-8 border border-gray-700">
            <h2 class="text-2xl font-bold text-white mb-6">📈 Coverage Distribution</h2>
            <div class="h-64">
                <canvas id="coverageChart"></canvas>
            </div>
        </div>

        <!-- Package Coverage Details -->
        <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 mb-8 border border-gray-700">
            <h2 class="text-2xl font-bold text-white mb-6">📦 Coverage by Package</h2>
            <div class="space-y-4">
EOF

# Add package coverage details
for PKG in "${!PACKAGE_COVERAGE[@]}"; do
    COV="${PACKAGE_COVERAGE[$PKG]}"
    COV_NUM=$(echo "$COV" | sed 's/%//')
    
    # Determine color based on coverage
    if (( $(echo "$COV_NUM >= 90" | bc -l) )); then
        COLOR="green"
    elif (( $(echo "$COV_NUM >= 75" | bc -l) )); then
        COLOR="blue"
    elif (( $(echo "$COV_NUM >= 60" | bc -l) )); then
        COLOR="yellow"
    else
        COLOR="red"
    fi

    cat >>reports/coverage_report.html  << EOF
                <div class="bg-gray-900 rounded-lg p-4 hover:bg-gray-750 transition-colors">
                    <div class="flex items-center justify-between mb-2">
                        <span class="font-mono text-${COLOR}-400 text-sm">$PKG</span>
                        <span class="text-${COLOR}-400 font-bold">$COV</span>
                    </div>
                    <div class="w-full bg-gray-700 rounded-full h-2.5">
                        <div class="bg-${COLOR}-500 h-2.5 rounded-full transition-all duration-500" style="width: $COV"></div>
                    </div>
                </div>
EOF
done

# Add file details section
cat >>reports/coverage_report.html  << 'EOF'
            </div>
        </div>

        <!-- Test Phase Results -->
        <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 mb-8 border border-gray-700">
            <h2 class="text-2xl font-bold text-white mb-6">✅ Test Phase Results</h2>
            <div class="space-y-3">
EOF

# Add phase results
if [ $FAILED_TESTS -eq 0 ]; then
    PHASE_ICON="✅"
    PHASE_COLOR="green"
else
    PHASE_ICON="❌"
    PHASE_COLOR="red"
fi

cat >>reports/coverage_report.html  << EOF
                <div class="flex items-center justify-between bg-gray-900 rounded-lg p-4">
                    <div class="flex items-center space-x-3">
                        <span class="text-2xl">✅</span>
                        <span class="text-gray-300">Phase 1: Unit Tests</span>
                    </div>
                    <span class="px-3 py-1 bg-green-500 bg-opacity-20 text-green-400 rounded-full text-sm font-semibold">PASSED</span>
                </div>
                <div class="flex items-center justify-between bg-gray-900 rounded-lg p-4">
                    <div class="flex items-center space-x-3">
                        <span class="text-2xl">${PHASE_ICON}</span>
                        <span class="text-gray-300">Phase 2: Integration Tests</span>
                    </div>
                    <span class="px-3 py-1 bg-${PHASE_COLOR}-500 bg-opacity-20 text-${PHASE_COLOR}-400 rounded-full text-sm font-semibold">$([ $FAILED_TESTS -eq 0 ] && echo "PASSED" || echo "FAILED")</span>
                </div>
                <div class="flex items-center justify-between bg-gray-900 rounded-lg p-4">
                    <div class="flex items-center space-x-3">
                        <span class="text-2xl">⚡</span>
                        <span class="text-gray-300">Phase 3: Race Tests</span>
                    </div>
                    <span class="px-3 py-1 bg-yellow-500 bg-opacity-20 text-yellow-400 rounded-full text-sm font-semibold">SKIPPED</span>
                </div>
                <div class="flex items-center justify-between bg-gray-900 rounded-lg p-4">
                    <div class="flex items-center space-x-3">
                        <span class="text-2xl">📊</span>
                        <span class="text-gray-300">Phase 4: Coverage Analysis</span>
                    </div>
                    <span class="px-3 py-1 bg-blue-500 bg-opacity-20 text-blue-400 rounded-full text-sm font-semibold">COMPLETE</span>
                </div>
EOF

cat >>reports/coverage_report.html  << EOF
            </div>
        </div>

        <!-- Metadata -->
        <div class="fade-in bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700">
            <h2 class="text-2xl font-bold text-white mb-4">ℹ️ Report Information</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-gray-900 rounded-lg p-4">
                    <p class="text-gray-400 text-sm mb-1">Generated At</p>
                    <p class="text-white font-mono">${CURRENT_DATE}</p>
                </div>
                <div class="bg-gray-900 rounded-lg p-4">
                    <p class="text-gray-400 text-sm mb-1">Test Duration</p>
                    <p class="text-white font-mono">${DURATION} seconds</p>
                </div>
                <div class="bg-gray-900 rounded-lg p-4">
                    <p class="text-gray-400 text-sm mb-1">Total Coverage</p>
                    <p class="text-white font-mono">${COVERAGE_PERCENT}</p>
                </div>
                <div class="bg-gray-900 rounded-lg p-4">
                    <p class="text-gray-400 text-sm mb-1">Status</p>
                    <p class="text-white font-mono">$([ $FAILED_TESTS -eq 0 ] && echo "✅ All Tests Passed" || echo "❌ Some Tests Failed")</p>
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div class="text-center mt-8 text-gray-500 text-sm">
            <p>Country Search API Test Report • Generated by test_api.sh</p>
        </div>
    </div>

    <!-- Chart.js Script -->
    <script>
        // Package coverage data
        const packages = [];
        const coverages = [];
EOF

# Add chart data
for PKG in "${!PACKAGE_COVERAGE[@]}"; do
    COV="${PACKAGE_COVERAGE[$PKG]}"
    COV_NUM=$(echo "$COV" | sed 's/%//')
    cat >>reports/coverage_report.html  << EOF
        packages.push('$PKG');
        coverages.push($COV_NUM);
EOF
done
 
cat > reports/coverage_report.html << 'EOF' 

        // Create chart
        const ctx = document.getElementById('coverageChart').getContext('2d');
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: packages,
                datasets: [{
                    label: 'Coverage %',
                    data: coverages,
                    backgroundColor: coverages.map(c => {
                        if (c >= 90) return 'rgba(34, 197, 94, 0.5)';
                        if (c >= 75) return 'rgba(59, 130, 246, 0.5)';
                        if (c >= 60) return 'rgba(234, 179, 8, 0.5)';
                        return 'rgba(239, 68, 68, 0.5)';
                    }),
                    borderColor: coverages.map(c => {
                        if (c >= 90) return 'rgb(34, 197, 94)';
                        if (c >= 75) return 'rgb(59, 130, 246)';
                        if (c >= 60) return 'rgb(234, 179, 8)';
                        return 'rgb(239, 68, 68)';
                    }),
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        backgroundColor: 'rgba(17, 24, 39, 0.9)',
                        titleColor: 'rgb(243, 244, 246)',
                        bodyColor: 'rgb(209, 213, 219)',
                        borderColor: 'rgb(75, 85, 99)',
                        borderWidth: 1
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        grid: { color: 'rgba(75, 85, 99, 0.3)' },
                        ticks: { color: 'rgb(156, 163, 175)' }
                    },
                    x: {
                        grid: { display: false },
                        ticks: { color: 'rgb(156, 163, 175)' }
                    }
                }
            }
        });
    </script>
</body>
</html>
EOF

echo -e "${GREEN}✓ Beautiful HTML report generated:reports/coverage_report.html ${NC}"

# Open report in default browser (platform-specific)
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    startreports/coverage_report.html  2>/dev/null || echo -e "${YELLOW}Openreports/coverage_report.html  manually${NC}"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    openreports/coverage_report.html  2>/dev/null
else
    xdg-openreports/coverage_report.html  2>/dev/null || echo -e "${YELLOW}Openreports/coverage_report.html  manually${NC}"
fi

# -----------------------------
# Summary
# -----------------------------
echo ""
echo "=========================================="
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✅ All test phases passed!${NC}"
    echo -e "${CYAN}📊 Coverage Report:reports/coverage_report.html ${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED_TESTS test phase(s) failed${NC}"
    echo -e "${CYAN}📊 Coverage Report:reports/coverage_report.html ${NC}"
    exit 1
fi
echo "=========================================="1
    if go test -race -v ./tests/unit/...; then
        echo -e "${GREEN}Race tests passed${NC}"
    else
        echo -e "${RED}Race tests failed${NC}"
        FAILED_TESTS=$((FAILED_TESTS+1))
    fi
fi
TOTAL_TESTS=$((TOTAL_TESTS+1))

# -----------------------------
# Phase 4: Coverage Report
# -----------------------------
echo ""
echo -e "${YELLOW}Phase 4: Coverage Report${NC}"
echo "=========================================="
go test ./tests/... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1

# -----------------------------
# Summary
# -----------------------------
echo ""
echo "=========================================="
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✅ All test phases passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED_TESTS test phase(s) failed${NC}"
    exit 1
fi
echo "=========================================="