
// ========================================
// FILE: internal/handler/health_handler.go
// ========================================
package handler

import (
	"net/http"
	
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/internal/service"
)

type HealthHandler struct {
	service service.CountryService
	logger  *logger.Logger
}

func NewHealthHandler(service service.CountryService, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	})
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	stats := h.service.GetCacheStats()
	
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"cache": map[string]int64{
			"hits":      stats.Hits,
			"misses":    stats.Misses,
			"size":      stats.Size,
			"evictions": stats.Evictions,
		},
	})
}
