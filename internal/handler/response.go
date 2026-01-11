
// ========================================
// FILE: internal/handler/response.go
// ========================================
package handler

import (
	"encoding/json"
	"net/http"
	
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	//"github.com/krishnajain872/country-search-api-cache/pkg/types"
	appErr "github.com/krishnajain872/country-search-api-cache/pkg/error"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, log *logger.Logger, err error) {
	appError := appErr.Get(err)
	
	log.Error("Request error",
		logger.String("type", string(appError.Type)),
		logger.String("message", appError.Message),
		logger.Int("status_code", appError.StatusCode))
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.StatusCode)
	
	response := Response{
		Success: false,
		Error: &ErrorResponse{
			Type:    string(appError.Type),
			Message: appError.Message,
			Context: appError.Context,
		},
	}
	
	json.NewEncoder(w).Encode(response)
}