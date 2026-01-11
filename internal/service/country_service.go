package service

import (
	"context"
	"time"
	
	"github.com/krishnajain872/country-search-api-cache/internal/domain"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/internal/repository"
	"github.com/krishnajain872/country-search-api-cache/pkg/cache"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

type CountryService interface {
	SearchCountry(ctx context.Context, name string) (*types.Country, error)
	GetCacheStats() types.Stats
}

type countryService struct {
	repo     repository.CountryRepository
	cache    cache.Cache
	logger   *logger.Logger
	cacheTTL time.Duration
}

func NewCountryService(
	repo repository.CountryRepository,
	cache cache.Cache,
	logger *logger.Logger,
	cacheTTL time.Duration,
) CountryService {
	return &countryService{
		repo:     repo,
		cache:    cache,
		logger:   logger,
		cacheTTL: cacheTTL,
	}
}

func (s *countryService) SearchCountry(ctx context.Context, name string) (*types.Country, error) {
	// Validate input
	searchReq, err := domain.NewSearchRequest(name)
	if err != nil {
		s.logger.Warn("Invalid search request",
			logger.String("name", name),
			logger.String("error", err.Error()))
		return nil, err
	}
	
	cacheKey := domain.CacheKey(searchReq)
	
	// Check cache
	s.logger.Debug("Checking cache",
		logger.String("cache_key", cacheKey))
	
	cachedValue, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		s.logger.Info("Cache hit",
			logger.String("country", name))
		
		if country, ok := cachedValue.(*types.Country); ok {
			return country, nil
		}
	}
	
	// Cache miss - fetch from API
	s.logger.Info("Cache miss, fetching from API",
		logger.String("country", name))
	
	country, err := s.repo.FindByName(ctx, searchReq.Name)
	if err != nil {
		s.logger.Error("Failed to fetch country",
			logger.String("country", name),
			logger.String("error", err.Error()))
		return nil, err
	}
	
	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, country, s.cacheTTL); err != nil {
		s.logger.Warn("Failed to cache country",
			logger.String("country", name),
			logger.String("error", err.Error()))
		// Don't return error, cache failure shouldn't fail the request
	}
	
	s.logger.Info("Successfully retrieved country",
		logger.String("country", country.Name))
	
	return country, nil
}

func (s *countryService) GetCacheStats() types.Stats {
	return s.cache.Stats()
}
