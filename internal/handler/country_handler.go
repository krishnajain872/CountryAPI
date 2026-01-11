
// ========================================
// FILE: internal/handler/country_handler.go
// ========================================
package handler

import (
	"net/http"
	
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/internal/service"
	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

type CountryHandler struct {
	service service.CountryService
	logger  *logger.Logger
}

func NewCountryHandler(service service.CountryService, logger *logger.Logger) *CountryHandler {
	return &CountryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CountryHandler) SearchCountry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, h.logger, appErr.Validation("method", "not allowed"))
		return
	}
	
	name := r.URL.Query().Get("name")
	if name == "" {
		err := appErr.Validation("name", "query parameter is required")
		WriteError(w, h.logger, err)
		return
	}
	
	h.logger.Info("Searching for country",
		logger.String("name", name))
	
	country, err := h.service.SearchCountry(r.Context(), name)
	if err != nil {
		WriteError(w, h.logger, err)
		return
	}
	
	WriteJSON(w, http.StatusOK, country)
}

