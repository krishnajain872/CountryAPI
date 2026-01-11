package unit
import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

// Mock service
type mockService struct {
	searchCountryFunc func(ctx context.Context, name string) (*types.Country, error)
	getCacheStatsFunc func() types.Stats
}

func (m *mockService) SearchCountry(ctx context.Context, name string) (*types.Country, error) {
	return m.searchCountryFunc(ctx, name)
}

func (m *mockService) GetCacheStats() types.Stats {
	return m.getCacheStatsFunc()
}

func TestCountryHandler_SearchCountry_Success(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	
	mockSvc := &mockService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			return &types.Country{
				Name:       "India",
				Capital:    "New Delhi",
				Currency:   "₹",
				Population: 1380004385,
			}, nil
		},
	}
	
	handler := NewCountryHandler(mockSvc, log)
	
	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	w := httptest.NewRecorder()
	
	handler.SearchCountry(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCountryHandler_SearchCountry_MissingName(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	
	mockSvc := &mockService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Service should not be called when name is missing")
			return nil, nil
		},
	}
	
	handler := NewCountryHandler(mockSvc, log)
	
	req := httptest.NewRequest(http.MethodGet, "/api/countries/search", nil)
	w := httptest.NewRecorder()
	
	handler.SearchCountry(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCountryHandler_SearchCountry_WrongMethod(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	
	mockSvc := &mockService{
		searchCountryFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Service should not be called for wrong method")
			return nil, nil
		},
	}
	
	handler := NewCountryHandler(mockSvc, log)
	
	req := httptest.NewRequest(http.MethodPost, "/api/countries/search?name=India", nil)
	w := httptest.NewRecorder()
	
	handler.SearchCountry(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
