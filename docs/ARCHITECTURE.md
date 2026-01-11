# Country Search API - Architecture Design

## Overview
This document outlines the architecture for a production-ready Country Search API with advanced in-memory caching, comprehensive error handling, and structured logging.

## Architecture Principles

### 1. **Clean Architecture / Layered Architecture**
- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic and orchestration
- **Repository Layer**: External API communication
- **Cache Layer**: In-memory data storage
- **Domain Layer**: Core business models and rules

### 2. **Dependency Injection**
All dependencies are injected through constructors, enabling:
- Easy testing with mocks
- Loose coupling between layers
- Better maintainability

## Component Architecture

### Cache Layer (`internal/cache/`)

#### Interface Design
```go
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Clear(ctx context.Context) error
    Stats() CacheStats
}

type CacheStats struct {
    Hits        int64
    Misses      int64
    Size        int64
    Evictions   int64
}
```

#### Advanced Features
1. **Thread Safety**: RWMutex for concurrent access
2. **TTL Support**: Time-to-live for cache entries
3. **LRU Eviction**: Optional Least Recently Used eviction
4. **Metrics Collection**: Hit rate, miss rate, size tracking
5. **Context Support**: Cancellation and timeout support

#### Implementation Strategy
```
┌──────────────────────────────────────┐
│     MemoryCache                      │
│  ┌────────────────────────────────┐  │
│  │  sync.RWMutex (Thread Safety)  │  │
│  └────────────────────────────────┘  │
│  ┌────────────────────────────────┐  │
│  │  map[string]CacheEntry         │  │
│  │  - value: interface{}          │  │
│  │  - expiry: time.Time           │  │
│  │  - lastAccess: time.Time       │  │
│  └────────────────────────────────┘  │
│  ┌────────────────────────────────┐  │
│  │  Metrics & Stats               │  │
│  └────────────────────────────────┘  │
└──────────────────────────────────────┘
```

### Repository Layer (`internal/repository/`)

#### Responsibilities
- HTTP communication with REST Countries API
- Request/response mapping
- Network error handling
- Timeout management
- Retry logic with exponential backoff

#### HTTP Client Configuration
```go
type HTTPClient struct {
    client     *http.Client
    baseURL    string
    timeout    time.Duration
    maxRetries int
    logger     logger.Logger
}
```

### Service Layer (`internal/service/`)

#### Cache-Aside Pattern Implementation
```
┌─────────────────────────────────────────┐
│          API Request                    │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│    Check Cache (Get)                    │
└─────────────┬───────────────────────────┘
              │
       ┌──────┴──────┐
       │             │
    Found         Not Found
       │             │
       ▼             ▼
   Return      ┌─────────────────────┐
   Cached      │ Call External API   │
   Data        └──────────┬──────────┘
                          │
                          ▼
                 ┌─────────────────────┐
                 │ Store in Cache      │
                 └──────────┬──────────┘
                            │
                            ▼
                     Return Fresh Data
```

### Error Handling Strategy

#### Error Types (`internal/domain/errors.go`)
```go
type ErrorType string

const (
    ErrorTypeValidation  ErrorType = "VALIDATION_ERROR"
    ErrorTypeNotFound    ErrorType = "NOT_FOUND"
    ErrorTypeExternal    ErrorType = "EXTERNAL_API_ERROR"
    ErrorTypeCache       ErrorType = "CACHE_ERROR"
    ErrorTypeInternal    ErrorType = "INTERNAL_ERROR"
    ErrorTypeTimeout     ErrorType = "TIMEOUT_ERROR"
)

type AppError struct {
    Type       ErrorType
    Message    string
    StatusCode int
    Err        error
    Context    map[string]interface{}
}
```

#### Error Handling Flow
```
Request → Handler → Service → Repository
   │         │         │          │
   └─────────┴─────────┴──────────┘
                  │
          Error Middleware
                  │
           ┌──────┴──────┐
           │             │
      Log Error    Format Response
           │             │
           └──────┬──────┘
                  │
           Return to Client
```

### Logging Strategy (`internal/logger/`)

#### Structured Logging Levels
- **DEBUG**: Cache operations, detailed flow
- **INFO**: API requests, successful operations
- **WARN**: Cache misses, retry attempts
- **ERROR**: Failed API calls, validation errors
- **FATAL**: Startup failures, critical errors

#### Log Context Enrichment
```go
type LogContext struct {
    RequestID  string
    Method     string
    Path       string
    Duration   time.Duration
    StatusCode int
    Error      error
}
```

#### Example Log Format (JSON)
```json
{
  "timestamp": "2026-01-10T10:30:45Z",
  "level": "INFO",
  "message": "Country search request",
  "request_id": "abc-123-def",
  "country_name": "India",
  "cache_hit": true,
  "duration_ms": 2.5
}
```

## Request Flow

### Complete Request Lifecycle
```
1. Client Request
   └─> GET /api/countries/search?name=India
       │
2. Middleware Chain
   ├─> Request ID Generation
   ├─> Logging (Request Start)
   ├─> Recovery (Panic Handler)
   └─> Timeout Context
       │
3. Handler Layer
   ├─> Parse & Validate Query Parameters
   ├─> Generate Cache Key
   └─> Call Service Layer
       │
4. Service Layer
   ├─> Check Cache (with context)
   │   ├─> Cache Hit → Return Cached Data
   │   └─> Cache Miss → Continue
   ├─> Call Repository Layer
   ├─> Transform Data to Domain Model
   ├─> Store in Cache (async or sync)
   └─> Return Data
       │
5. Handler Layer
   ├─> Format Response
   ├─> Set Headers
   └─> Send JSON Response
       │
6. Middleware Chain
   └─> Logging (Request Complete)
```

## Configuration Management

### Environment-Based Configuration
```go
type Config struct {
    Server ServerConfig
    Cache  CacheConfig
    API    APIConfig
    Logger LoggerConfig
}

type CacheConfig struct {
    TTL            time.Duration
    MaxSize        int
    EvictionPolicy string
    CleanupInterval time.Duration
}
```

### Configuration Sources
1. Environment variables (12-factor app)
2. Configuration files (YAML/JSON)
3. Defaults with override capability

## Testing Strategy

### Unit Tests
- **Cache Layer**: Concurrent access, TTL, eviction
- **Service Layer**: Business logic with mocked dependencies
- **Handler Layer**: HTTP request/response with mocked service
- **Repository Layer**: HTTP client with mock server

### Race Condition Tests
```bash
go test -race ./internal/cache/...
go test -race ./internal/service/...
```

### Integration Tests
- Full request flow
- Cache integration with service
- External API interaction

### Test Coverage Target
- Minimum 80% code coverage
- 100% coverage for critical paths (cache, service logic)

## Deployment Considerations

### Graceful Shutdown
```go
func gracefulShutdown(server *http.Server, timeout time.Duration) {
    // Wait for interrupt signal
    // Stop accepting new requests
    // Complete in-flight requests
    // Close resources (cache, connections)
    // Exit
}
```

### Health Checks
- `/health/live`: Liveness probe
- `/health/ready`: Readiness probe (checks cache, API availability)

### Metrics Endpoints
- `/metrics`: Prometheus-compatible metrics
  - Request count, duration
  - Cache hit/miss rates
  - Error rates

## Security Considerations

1. **Input Validation**: Sanitize all user inputs
2. **Rate Limiting**: Prevent abuse (optional middleware)
3. **Timeout Management**: Prevent resource exhaustion
4. **Error Information**: Don't leak sensitive data in errors

## Performance Optimization

1. **Connection Pooling**: Reuse HTTP connections
2. **Context Propagation**: Enable cancellation
3. **Async Cache Updates**: Non-blocking cache writes
4. **Response Compression**: GZIP middleware
5. **Cache Prewarming**: Popular countries on startup

## Scalability Considerations

### Current Implementation (Single Instance)
- In-memory cache per instance
- Suitable for single-server deployment

### Future Enhancements
- Distributed cache (Redis)
- Horizontal scaling with shared cache
- Cache invalidation strategies

## Monitoring & Observability

### Key Metrics to Track
- Request throughput (req/sec)
- Response latency (p50, p95, p99)
- Cache hit ratio
- External API error rate
- Memory usage

### Logging Strategy
- Structured JSON logs
- Request ID tracking
- Error stack traces
- Performance metrics in logs

## Technology Stack

- **Language**: Go 1.21+
- **HTTP Router**: Standard library `net/http` or Chi/Gorilla
- **Logging**: Zap (uber-go/zap)
- **Testing**: Standard library + testify
- **Documentation**: OpenAPI 3.0

## Summary

This architecture provides:
- ✅ Clean separation of concerns
- ✅ Thread-safe concurrent operations
- ✅ Comprehensive error handling
- ✅ Structured logging with context
- ✅ Testable components
- ✅ Production-ready patterns
- ✅ Scalability path
- ✅ Observability built-in