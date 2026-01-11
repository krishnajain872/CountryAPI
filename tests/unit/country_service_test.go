package unit

import (
	"context"
	"errors"
	"testing"
	"time"
	
	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"github.com/krishnajain872/country-search-api-cache/internal/service"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/pkg/cache"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

// Mock repository
type mockRepository struct {
	findByNameFunc func(ctx context.Context, name string) (*types.Country, error)
}

func (m *mockRepository) FindByName(ctx context.Context, name string) (*types.Country, error) {
	return m.findByNameFunc(ctx, name)
}

func TestCountryService_SearchCountry_CacheHit(t *testing.T) {
	// Setup
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	memCache := cache.NewMemoryCache(5*time.Minute, 100)
	
	mockRepo := &mockRepository{
		findByNameFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Repository should not be called on cache hit")
			return nil, errors.New("should not be called")
		},
	}
	
	service := service.NewCountryService(mockRepo, memCache, log, 5*time.Minute)
	
	// Pre-populate cache
	country := &types.Country{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "₹",
		Population: 1380004385,
	}
	memCache.Set(context.Background(), "country:india", country, 0)
	
	// Test
	result, err := service.SearchCountry(context.Background(), "India")
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.Name != "India" {
		t.Errorf("Expected India, got %s", result.Name)
	}
	
	// Verify cache hit
	stats := service.GetCacheStats()
	if stats.Hits < 1 {
		t.Error("Expected at least one cache hit")
	}
}

func TestCountryService_SearchCountry_CacheMiss(t *testing.T) {
	// Setup
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	memCache := cache.NewMemoryCache(5*time.Minute, 100)
	
	repoCallCount := 0
	mockRepo := &mockRepository{
		findByNameFunc: func(ctx context.Context, name string) (*types.Country, error) {
			repoCallCount++
			return &types.Country{
				Name:       "India",
				Capital:    "New Delhi",
				Currency:   "₹",
				Population: 1380004385,
			}, nil
		},
	}
	
	service := service.NewCountryService(mockRepo, memCache, log, 5*time.Minute)
	
	// Test
	result, err := service.SearchCountry(context.Background(), "India")
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.Name != "India" {
		t.Errorf("Expected India, got %s", result.Name)
	}
	
	if repoCallCount != 1 {
		t.Errorf("Expected repository to be called once, got %d", repoCallCount)
	}
	
	// Verify cache miss then population
	stats := service.GetCacheStats()
	if stats.Misses < 1 {
		t.Error("Expected at least one cache miss")
	}
}

func TestCountryService_SearchCountry_InvalidInput(t *testing.T) {
	cfg := &config.LoggerConfig{
		Environment: config.Development,
		Mode:        []config.LogModeType{config.ConsoleMode},
		Severity:    config.InfoSeverity,
	}
	log := logger.NewLogger(cfg)
	memCache := cache.NewMemoryCache(5*time.Minute, 100)
	
	mockRepo := &mockRepository{
		findByNameFunc: func(ctx context.Context, name string) (*types.Country, error) {
			t.Error("Repository should not be called for invalid input")
			return nil, errors.New("should not be called")
		},
	}
	
	service := service.NewCountryService(mockRepo, memCache, log, 5*time.Minute)
	
	// Test with empty name
	_, err := service.SearchCountry(context.Background(), "")
	
	if err == nil {
		t.Error("Expected validation error for empty name")
	}
}
