// internal/config/loader.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfig loads config from environment or `.env` file and returns Config
func LoadConfig(envFile string) *Config {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("⚠️ Warning: Could not load env file: %v", err)
		}
	}

	cfg := &Config{
		Server:   loadServerConfig(),
		Cache:    loadCacheConfig(),
		External: loadExternalAPIConfig(),
		Logger:   loadLoggerConfig(),
	}

	Store(cfg) // save singleton
	return cfg
}

// ====================
// Section loaders
// ====================

func loadServerConfig() *ServerConfig {
	cfg := defaultServerConfig()
	cfg.Host = getEnv("SERVER_HOST", cfg.Host)
	cfg.Port = getEnv("SERVER_PORT", cfg.Port)
	cfg.ReadTimeout = getEnvAsDuration("SERVER_READ_TIMEOUT", cfg.ReadTimeout)
	cfg.WriteTimeout = getEnvAsDuration("SERVER_WRITE_TIMEOUT", cfg.WriteTimeout)
	return cfg
}

func loadCacheConfig() *CacheConfig {
	cfg := defaultCacheConfig()
	cfg.TTL = getEnvAsDuration("CACHE_TTL", cfg.TTL)
	cfg.MaxSize = getEnvAsInt("CACHE_MAX_SIZE", cfg.MaxSize)
	return cfg
}

func loadExternalAPIConfig() *ExternalAPIConfig {
	cfg := defaultExternalAPIConfig()
	cfg.RestCountriesURL = getEnv("REST_COUNTRIES_API_URL", cfg.RestCountriesURL)
	cfg.Timeout = getEnvAsDuration("API_TIMEOUT", cfg.Timeout)
	return cfg
}

func loadLoggerConfig() *LoggerConfig {
	cfg := defaultLoggerConfig()
	cfg.Level = getEnv("LOG_LEVEL", cfg.Level)
	cfg.Format = getEnv("LOG_FORMAT", cfg.Format)
	return cfg
}

// ====================
// Environment helpers
// ====================

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := time.ParseDuration(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
