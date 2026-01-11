
// ========================================
// FILE: cmd/server/main.go (UPDATED)
// ========================================
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/krishnajain872/country-search-api-cache/internal/config"
	"github.com/krishnajain872/country-search-api-cache/internal/handler"
	"github.com/krishnajain872/country-search-api-cache/internal/logger"
	"github.com/krishnajain872/country-search-api-cache/internal/repository"
	"github.com/krishnajain872/country-search-api-cache/internal/service"
	"github.com/krishnajain872/country-search-api-cache/pkg/cache"
	"github.com/krishnajain872/country-search-api-cache/pkg/httpclient"
)

func main() {
	// ----------------------------- 
	// 1️⃣ Load configuration
	// ----------------------------- 
	cfg := config.LoadConfig(".env")
	
	// Initialize logger with config
	log := logger.NewLogger(cfg.Logger)
	defer log.Sync()
	
	log.Info("Starting Country Search API",
		logger.String("host", cfg.Server.Host),
		logger.String("port", cfg.Server.Port))
	
	// ----------------------------- 
	// 2️⃣ Initialize cache
	// ----------------------------- 
	memCache := cache.NewMemoryCache(cfg.Cache.TTL, cfg.Cache.MaxSize)
	log.Info("Cache initialized",
		logger.Any("ttl", cfg.Cache.TTL),
		logger.Int("max_size", cfg.Cache.MaxSize))
	
	// ----------------------------- 
	// 3️⃣ Initialize HTTP client
	// ----------------------------- 
	httpClient := httpclient.NewClient(
		cfg.External.RestCountriesURL,
		cfg.External.Timeout,
	)
	
	// ----------------------------- 
	// 4️⃣ Initialize repository
	// ----------------------------- 
	repo := repository.NewRestCountriesRepo(httpClient, log)
	
	// ----------------------------- 
	// 5️⃣ Initialize service
	// ----------------------------- 
	countryService := service.NewCountryService(repo, memCache, log, cfg.Cache.TTL)
	
	// ----------------------------- 
	// 6️⃣ Initialize router
	// ----------------------------- 
	router := handler.NewRouter(countryService, log)
	
	// ----------------------------- 
	// 7️⃣ Create HTTP server
	// ----------------------------- 
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	
	// ----------------------------- 
	// 8️⃣ Start server in goroutine
	// ----------------------------- 
	go func() {
		log.Info("Server listening",
			logger.String("address", server.Addr))
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start",
				logger.String("error", err.Error()))
		}
	}()
	
	// ----------------------------- 
	// 9️⃣ Wait for interrupt signal
	// ----------------------------- 
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("Shutting down server...")
	
	// ----------------------------- 
	// 🔟 Graceful shutdown
	// ----------------------------- 
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown",
			logger.String("error", err.Error()))
	}
	
	log.Info("Server stopped gracefully")
}