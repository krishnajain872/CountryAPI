// cmd/server/main.go
package main

import (
	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"github.com/krishnajain872/country-search-api-cache/internal/domain"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
)

func main() {
	// -----------------------------
	// 1️⃣ Load configuration
	// -----------------------------
	cfg := config.LoadConfig(".env")

	// Initialize logger with config
	log := logger.NewLogger(cfg.Logger)

	// -----------------------------
	// 2️⃣ Server start log
	// -----------------------------
	log.Info("Server starting...",
		logger.String("host", cfg.Server.Host),
		logger.String("port", cfg.Server.Port),
	)

	// -----------------------------
	// 3️⃣ Create validated search request
	// -----------------------------
	searchReq, err := domain.NewSearchRequest("India")
	if err != nil {
		log_string := handleError(err)
		log.Error("Search request creation failed",
			logger.String("error", log_string),
		)
		return
	}

	log.Info("Search request created",
		logger.Any("request", searchReq),
	)

	// -----------------------------
	// 4️⃣ Create country entity
	// -----------------------------
	country := domain.NewCountry(
		"India",
		"New Delhi",
		"INR",
		1_428_000_000,
	)

	// -----------------------------
	// 5️⃣ Validate country entity
	// -----------------------------
	if err := domain.ValidateCountry(country); err != nil {
		log_string := handleError(err)
		log.Error("Country validation failed",
			logger.String("error", log_string),
		)
		return
	}

	log.Info("Country validated successfully",
		logger.Any("country", country),
	)
}
