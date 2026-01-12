// ========================================
// FILE: tests/unit/country_handler_test.go
// Handler layer unit tests
// ========================================
package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"github.com/krishnajain872/country-search-api-cache/internal/handler"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

// Mock service
type mockCountryService struct {
	searchCountryFunc func(ctx context.Context, name string) (*types.Country, error)
	getCacheStatsFunc func() types.Stats
}

func (m *mockCountryService) SearchCountry(ctx context.Context, name string) (*types.Country, error) {
	return m.searchCountryFunc(ctx, name)
}

func (m *mockCountryService) GetCacheStats() types.Stats {
	return m.getCacheStatsFunc()
}

func TestCountryHandler_SearchCountry_Success(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
		FilePath:    "logs/test.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
	log := logger.NewLogger(cfg)

	mockSvc := &mockCountryService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			return &types.Country{
				Name:       "India",
				Capital:    "New Delhi",
				Currency:   "₹",
				Population: 1380004385,
			}, nil
		},
	}

	h := handler.NewCountryHandler(mockSvc, log)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	w := httptest.NewRecorder()

	h.SearchCountry(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response handler.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success=true")
	}
}

func TestCountryHandler_SearchCountry_MissingName(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
		FilePath:    "logs/test.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
	log := logger.NewLogger(cfg)

	mockSvc := &mockCountryService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Service should not be called when name is missing")
			return nil, nil
		},
	}

	h := handler.NewCountryHandler(mockSvc, log)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search", nil)
	w := httptest.NewRecorder()

	h.SearchCountry(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response handler.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false")
	}

	if response.Error == nil {
		t.Error("Expected error in response")
	}
}

func TestCountryHandler_SearchCountry_WrongMethod(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
		FilePath:    "logs/test.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
	log := logger.NewLogger(cfg)

	mockSvc := &mockCountryService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Service should not be called for wrong method")
			return nil, nil
		},
	}

	h := handler.NewCountryHandler(mockSvc, log)

	req := httptest.NewRequest(http.MethodPost, "/api/countries/search?name=India", nil)
	w := httptest.NewRecorder()

	h.SearchCountry(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCountryHandler_SearchCountry_NotFound(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
		FilePath:    "logs/test.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
	log := logger.NewLogger(cfg)

	mockSvc := &mockCountryService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			return nil, appErr.ResourceNotFound("country")
		},
	}

	h := handler.NewCountryHandler(mockSvc, log)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=InvalidCountry", nil)
	w := httptest.NewRecorder()

	h.SearchCountry(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response handler.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false")
	}

	if response.Error.Type != "NOT_FOUND" {
		t.Errorf("Expected NOT_FOUND error, got %s", response.Error.Type)
	}
}

func TestCountryHandler_SearchCountry_ValidationError(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
		FilePath:    "logs/test.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
	log := logger.NewLogger(cfg)

	mockSvc := &mockCountryService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			return nil, appErr.Validation("name", "must be at least 2 characters")
		},
	}

	h := handler.NewCountryHandler(mockSvc, log)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=A", nil)
	w := httptest.NewRecorder()

	h.SearchCountry(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response handler.Response
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success=false")
	}

	if response.Error.Type != "VALIDATION_ERROR" {
		t.Errorf("Expected VALIDATION_ERROR, got %s", response.Error.Type)
	}
}
