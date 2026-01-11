// ========================================
// FILE: internal/handler/router.go
// ========================================
package handler

import (
	"net/http"
	
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/internal/service"
)

func NewRouter(
	countryService service.CountryService,
	logger *logger.Logger,
) http.Handler {
	mux := http.NewServeMux()
	
	// Handlers
	countryHandler := NewCountryHandler(countryService, logger)
	healthHandler := NewHealthHandler(countryService, logger)
	
	// Routes
	mux.HandleFunc("/api/countries/search", countryHandler.SearchCountry)
	mux.HandleFunc("/health/live", healthHandler.Live)
	mux.HandleFunc("/health/ready", healthHandler.Ready)
	mux.HandleFunc("/metrics", healthHandler.Metrics)
	
	// Apply middleware
	handler := CORSMiddleware(mux)
	handler = RequestIDMiddleware(handler)
	handler = LoggingMiddleware(logger)(handler)
	handler = RecoveryMiddleware(logger)(handler)
	
	return handler
}
