# 🚀 Country Search API - Complete Guide

## 📋 Table of Contents
1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Installation & Setup](#installation--setup)
4. [Running the Application](#running-the-application)
5. [Testing Guide](#testing-guide)
6. [API Documentation](#api-documentation)
7. [Script Usage Guide](#script-usage-guide)
8. [Troubleshooting](#troubleshooting)

---

## 🎯 Project Overview

**Country Search API** is a high-performance REST API that searches for country information with intelligent caching.

### ✨ Key Features
- 🔍 Search countries by name
- ⚡ In-memory caching (30x faster on cache hits!)
- 🔒 Thread-safe concurrent operations
- 📊 Performance metrics and monitoring
- 🏥 Health check endpoints
- 🔐 URL encoding for special characters

### 🏗️ Architecture
```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│      Handler Layer              │
│  (HTTP Request/Response)        │
└──────────┬──────────────────────┘
           │
           ▼
┌─────────────────────────────────┐
│      Service Layer              │
│  (Business Logic + Cache)       │
└──────────┬──────────────────────┘
           │
    ┌──────┴──────┐
    ▼             ▼
┌────────┐   ┌──────────┐
│ Cache  │   │Repository│
└────────┘   └─────┬────┘
                   │
                   ▼
           ┌───────────────┐
           │ External API  │
           │(restcountries)│
           └───────────────┘
```

---

## 📦 Prerequisites

### Required Software
1. **Go (Golang)** - Version 1.21 or higher
   - Download: https://go.dev/dl/
   - Verify: `go version`

2. **Git** (for cloning the repository)
   - Download: https://git-scm.com/downloads
   - Verify: `git --version`

3. **Code Editor** (Optional but recommended)
   - VS Code: https://code.visualstudio.com/
   - GoLand: https://www.jetbrains.com/go/

### Optional Tools
- **Postman** - For API testing
  - Download: https://www.postman.com/downloads/
- **golangci-lint** - For code quality
  - Install: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **goimports** - For code formatting
  - Install: `go install golang.org/x/tools/cmd/goimports@latest`

---

## 🛠️ Installation & Setup

### Step 1: Clone the Repository
```bash
# Using HTTPS
git clone https://github.com/your-username/country-search-api.git

# OR using the existing folder
cd C:\Users\krish\Downloads\countery_search_api_cache_go_lang_assigememt
```

### Step 2: Navigate to Project Directory
```bash
cd countery_search_api_cache_go_lang_assigememt
```

### Step 3: Download Dependencies
```bash
go mod download
go mod tidy
```

### Step 4: Verify Setup
```bash
# Check if all packages are downloaded
go list -m all

# Verify Go environment
go env GOPATH
go env GOROOT
```

### Step 5: Configure Environment (Optional)
The `.env` file is already configured with defaults. You can modify it if needed:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8000

# Cache Configuration
CACHE_TTL=5m
CACHE_MAX_SIZE=1000

# External API
REST_COUNTRIES_API_URL=https://restcountries.com/v3.1
API_TIMEOUT=10s

# Logging
LOG_LEVEL=info
LOG_MODE=console,file
LOG_FILE_PATH=logs/app.log
```

---

## 🏃 Running the Application

### Method 1: Using Go Run (Quick Start)
```bash
# Run directly
go run cmd/server/main.go

# OR using the script (Linux/Mac/Git Bash)
chmod +x scripts/build.sh
./scripts/build.sh run
```

### Method 2: Build Then Run
```bash
# Build the binary
./scripts/build.sh build

# Run the binary
./bin/server
```

### Method 3: Using Make (if available)
```bash
make run
```

### Verify Server is Running
You should see:
```
INFO    Starting Country Search API     {"host": "localhost", "port": "8000"}
INFO    Cache initialized               {"ttl": "5m", "max_size": 1000}
INFO    Server listening                {"address": "localhost:8000"}
```

**✅ Server is ready at:** `http://localhost:8000`

---

## 🧪 Testing Guide

### Quick Test (All Tests)
```bash
# Run all tests
./scripts/build.sh test

# OR
go test ./... -v
```

### Comprehensive Test Suite
```bash
# Run the complete test suite with phases
./scripts/build.sh test-all
```

**Expected Output:**
```
==========================================
Running Comprehensive Test Suite
==========================================

📦 Phase 1: Unit Tests
==========================================
✅ Unit tests passed

💾 Phase 2: Cache Layer Tests
==========================================
✅ Cache tests passed

⚙️  Phase 3: Service Layer Tests
==========================================
✅ Service tests passed

🌐 Phase 4: Handler Layer Tests
==========================================
✅ Handler tests passed

🔗 Phase 5: Integration Tests
==========================================
✅ Integration tests passed

==========================================
Test Summary
==========================================
Total Phases: 5
Passed: 5
Failed: 0

✅ All test phases passed! 🎉
```

### Running Specific Test Suites

#### 1. Unit Tests Only
```bash
./scripts/build.sh test-unit
```

#### 2. Integration Tests Only
```bash
./scripts/build.sh test-integration
```

#### 3. Race Condition Tests
```bash
# Test for race conditions (concurrency bugs)
./scripts/build.sh race

# OR test all layers
./scripts/build.sh race-all
```

**What is Race Detection?**
- Detects concurrent access bugs
- Ensures thread-safe operations
- Critical for production readiness

### Coverage Reports

#### Generate Full Coverage
```bash
./scripts/build.sh coverage
```

**What You Get:**
- ✅ Coverage summary by package
- ✅ Total coverage percentage
- ✅ HTML report: `coverage.html`
- ✅ Text summary: `coverage_summary.txt`

**View HTML Report:**
```bash
# Windows
start coverage.html

# Mac
open coverage.html

# Linux
xdg-open coverage.html
```

#### Coverage by Package
```bash
./scripts/build.sh coverage-pkg
```

**Example Output:**
```
📦 Cache Layer:
pkg/cache/memory_cache.go:25:    NewMemoryCache    90.0%
pkg/cache/memory_cache.go:40:    Get               95.5%
pkg/cache/memory_cache.go:60:    Set               92.3%

⚙️  Service Layer:
internal/service/country_service.go:30:  SearchCountry  88.2%

🌐 Handler Layer:
internal/handler/country_handler.go:20:  SearchCountry  85.7%

Total Coverage: 87.3%
```

### Understanding Test Output

#### ✅ Successful Test
```
=== RUN   TestMemoryCache_Get
--- PASS: TestMemoryCache_Get (0.00s)
PASS
ok      pkg/cache       0.234s
```

#### ❌ Failed Test
```
=== RUN   TestMemoryCache_Get
    cache_test.go:45: Expected 'value', got 'wrong-value'
--- FAIL: TestMemoryCache_Get (0.00s)
FAIL
```

---

## 📡 API Documentation

### Base URL
```
http://localhost:8000
```

### Endpoints

#### 1. Search Country
**Endpoint:** `GET /api/countries/search`

**Query Parameters:**
| Parameter | Type   | Required | Description           |
|-----------|--------|----------|-----------------------|
| name      | string | Yes      | Country name to search|

**Example Requests:**
```bash
# Simple country
curl "http://localhost:8000/api/countries/search?name=India"

# Country with spaces
curl "http://localhost:8000/api/countries/search?name=United%20States"

# Country with special characters
curl "http://localhost:8000/api/countries/search?name=C%C3%B4te+d%27Ivoire"
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "name": "India",
    "capital": "New Delhi",
    "currency": "₹",
    "population": 1428000000
  }
}
```

**Error Response (400 - Missing Parameter):**
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "name query parameter is required",
    "context": {
      "field": "name"
    }
  }
}
```

**Error Response (404 - Not Found):**
```json
{
  "success": false,
  "error": {
    "type": "NOT_FOUND",
    "message": "country not found"
  }
}
```

#### 2. Health Check (Liveness)
**Endpoint:** `GET /health/live`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "status": "alive"
  }
}
```

#### 3. Health Check (Readiness)
**Endpoint:** `GET /health/ready`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "status": "ready"
  }
}
```

#### 4. Metrics
**Endpoint:** `GET /metrics`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "cache": {
      "hits": 42,
      "misses": 8,
      "size": 15,
      "evictions": 2
    }
  }
}
```

**Metrics Explanation:**
- **hits**: Number of cache hits (data found in cache)
- **misses**: Number of cache misses (data fetched from API)
- **size**: Current number of items in cache
- **evictions**: Number of items removed due to TTL or size limits

---

## 🛠️ Script Usage Guide

### build.sh - Complete Reference

**Location:** `scripts/build.sh`

**Make Executable (First Time Only):**
```bash
chmod +x scripts/build.sh
```

### 📚 All Available Commands

#### 🏗️ Build & Run Commands

| Command | Description | Example |
|---------|-------------|---------|
| `build` | Build the binary | `./scripts/build.sh build` |
| `run` | Run the server | `./scripts/build.sh run` |
| `clean` | Remove build artifacts | `./scripts/build.sh clean` |

**Example:**
```bash
# Build project
./scripts/build.sh build
# Output: bin/server created

# Run server
./scripts/build.sh run
# Server starts on http://localhost:8000

# Clean everything
./scripts/build.sh clean
# Removes: bin/, coverage.out, coverage.html
```

#### 🧪 Testing Commands

| Command | Description | What It Tests |
|---------|-------------|---------------|
| `test` | Run all tests | All packages |
| `test-unit` | Unit tests only | tests/unit/... |
| `test-integration` | Integration tests | tests/integration/... |
| `test-all` | Full test suite | All with phases |
| `race` | Race detector (unit) | Concurrency bugs |
| `race-all` | Race detector (all) | All layers |

**Example:**
```bash
# Quick test
./scripts/build.sh test

# Comprehensive test suite
./scripts/build.sh test-all

# Check for race conditions
./scripts/build.sh race-all
```

#### 📊 Coverage Commands

| Command | Description | Output |
|---------|-------------|--------|
| `coverage` | Full coverage report | HTML + Summary |
| `coverage-unit` | Unit test coverage | HTML |
| `coverage-pkg` | Package breakdown | Console |

**Example:**
```bash
# Generate full coverage
./scripts/build.sh coverage

# View by package
./scripts/build.sh coverage-pkg
```

**Output Files:**
- `coverage.out` - Raw coverage data
- `coverage.html` - Visual coverage report
- `coverage_summary.txt` - Text summary

#### 🎨 Code Quality Commands

| Command | Description | Requires |
|---------|-------------|----------|
| `fmt` | Format code | goimports (optional) |
| `lint` | Run linter | golangci-lint |
| `vet` | Go vet checks | Built-in |

**Example:**
```bash
# Format code
./scripts/build.sh fmt

# Run linter
./scripts/build.sh lint

# Static analysis
./scripts/build.sh vet
```

#### 📦 Dependency Commands

| Command | Description | When to Use |
|---------|-------------|-------------|
| `deps` | Download & tidy | After clone |
| `deps-update` | Update all deps | Monthly |

**Example:**
```bash
# Initial setup
./scripts/build.sh deps

# Update dependencies
./scripts/build.sh deps-update
```

#### 🐳 Docker Commands

| Command | Description | Output |
|---------|-------------|--------|
| `docker-build` | Build image | country-search-api:latest |
| `docker-run` | Run container | Starts on port 8000 |

**Example:**
```bash
# Build Docker image
./scripts/build.sh docker-build

# Run in Docker
./scripts/build.sh docker-run
```

#### ℹ️ Utility Commands

| Command | Description |
|---------|-------------|
| `help` | Show all commands |
| `version` | Show Go version |

**Example:**
```bash
# Show help
./scripts/build.sh help

# Check Go version
./scripts/build.sh version
```

---

## 🔧 Troubleshooting

### Common Issues & Solutions

#### 1. "Go is not installed"
**Error:**
```
❌ Go is not installed. Please install Go to continue.
```

**Solution:**
1. Download Go from https://go.dev/dl/
2. Install for your OS
3. Verify: `go version`
4. Restart terminal

#### 2. "Permission denied" (Linux/Mac)
**Error:**
```
bash: ./scripts/build.sh: Permission denied
```

**Solution:**
```bash
chmod +x scripts/build.sh
```

#### 3. Tests Fail with "cannot find package"
**Error:**
```
cache_test.go:8:2: cannot find package "github.com/..."
```

**Solution:**
```bash
# Download dependencies
go mod download
go mod tidy

# Or use script
./scripts/build.sh deps
```

#### 4. Integration Tests Fail
**Error:**
```
Integration tests failed (may need internet)
```

**Solution:**
- Check internet connection
- REST Countries API might be down
- Run without integration:
  ```bash
  ./scripts/build.sh test-unit
  ```

#### 5. Race Detector Fails
**Error:**
```
WARNING: DATA RACE
```

**Solution:**
- This indicates a concurrency bug
- Check the specific file and line number
- Review mutex usage in that code

#### 6. Port Already in Use
**Error:**
```
bind: address already in use
```

**Solution:**
```bash
# Find process using port 8000
# Windows
netstat -ano | findstr :8000

# Linux/Mac
lsof -i :8000

# Kill the process
# Windows
taskkill /PID <PID> /F

# Linux/Mac
kill -9 <PID>

# OR change port in .env
SERVER_PORT=8001
```

#### 7. Coverage Report Not Opening
**Solution:**
```bash
# Manually open
# Windows
start coverage.html

# Mac
open coverage.html

# Linux
xdg-open coverage.html

# OR open in browser
# Navigate to: file:///path/to/coverage.html
```

---

## 📊 Testing Checklist

Before submitting or deploying:

### ✅ Pre-Deployment Checklist
```bash
# 1. Clean build
./scripts/build.sh clean
./scripts/build.sh build

# 2. Run all tests
./scripts/build.sh test-all

# 3. Check race conditions
./scripts/build.sh race-all

# 4. Generate coverage (aim for >85%)
./scripts/build.sh coverage

# 5. Format code
./scripts/build.sh fmt

# 6. Run linter (if installed)
./scripts/build.sh lint

# 7. Run vet
./scripts/build.sh vet
```

### ✅ Expected Results
- ✅ All test phases passed
- ✅ No race conditions detected
- ✅ Coverage > 85%
- ✅ No linting errors
- ✅ No vet warnings

---

## 🎓 Learning Path

### For Beginners

#### Week 1: Understand the Basics
1. Read `docs/ARCHITECTURE.md`
2. Understand the API endpoints
3. Make API calls with cURL
4. Import Postman collection (see next section)

#### Week 2: Explore the Code
1. Start with `cmd/server/main.go`
2. Follow the flow: Handler → Service → Repository
3. Understand caching logic in `pkg/cache/`
4. Read error handling in `pkg/error/`

#### Week 3: Testing
1. Run all tests and understand output
2. Read test files in `tests/unit/`
3. Try modifying tests
4. Generate and analyze coverage

#### Week 4: Advanced
1. Add a new endpoint
2. Implement additional cache strategies
3. Add more validation rules
4. Contribute improvements

---

## 📮 Postman Collection

See the separate **Postman Collection** document for:
- Pre-configured requests
- Environment variables
- Test scripts
- Example responses

---

## 🎉 Success Criteria

You know you've succeeded when:

1. ✅ Server starts without errors
2. ✅ All API endpoints return correct responses
3. ✅ All tests pass: `./scripts/build.sh test-all`
4. ✅ No race conditions: `./scripts/build.sh race-all`
5. ✅ Coverage > 85%: `./scripts/build.sh coverage`
6. ✅ Can search countries with special characters
7. ✅ Cache hits are faster than cache misses
8. ✅ Health checks return success
9. ✅ Metrics show accurate stats
10. ✅ You understand the architecture!

---

## 📞 Support

If you encounter issues:
1. Check this documentation first
2. Review error messages carefully
3. Check `logs/app.log` for details
4. Search for similar issues on GitHub
5. Ask for help with specific error messages

---

## 🚀 Quick Start Summary

```bash
# 1. Clone & Setup
git clone <repo-url>
cd countery_search_api_cache_go_lang_assigememt
go mod download

# 2. Run Server
./scripts/build.sh run

# 3. Test API
curl "http://localhost:8000/api/countries/search?name=India"

# 4. Run Tests
./scripts/build.sh test-all

# 5. Check Coverage
./scripts/build.sh coverage
```
### TEST LIVE SERVER 
https://countryapi-v152.onrender.com/health/live

### USE POSTMAT COLLECTION FROM ./tests/postman_collection


 
