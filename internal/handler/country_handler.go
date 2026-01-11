// ========================================
// FILE: internal/handler/country_handler.go
// ========================================
package handler

import (
	"net/http"
	"net/url"

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

	// Get the raw query parameter
	encodedName := r.URL.Query().Get("name")
	if encodedName == "" {
		err := appErr.Validation("name", "query parameter is required")
		WriteError(w, h.logger, err)
		return
	}

	//   URL decode the name parameter
	// This handles both already-decoded and encoded inputs
	name, err := url.QueryUnescape(encodedName)
	if err != nil {
		h.logger.Warn("Failed to decode URL parameter",
			logger.String("encoded_name", encodedName),
			logger.String("error", err.Error()))
		
		// If decoding fails, use the original value
		name = encodedName
	}

	h.logger.Info("Searching for country",
		logger.String("name", name),
		logger.String("encoded", encodedName))

	country, err := h.service.SearchCountry(r.Context(), name)
	if err != nil {
		WriteError(w, h.logger, err)
		return
	}

	WriteJSON(w, http.StatusOK, country)
}
