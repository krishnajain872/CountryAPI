// internal/config/loader.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"
	"strings"
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

// loadLoggerConfig loads logger config from env variables
func loadLoggerConfig() *LoggerConfig {
	cfg := defaultLoggerConfig()

	// Environment: development, staging, production
	env := getEnv("LOG_ENV", string(cfg.Environment))
	cfg.Environment = LogEnvType(env)

	// Severity level: debug, info, warn, error, critical
	sev := getEnv("LOG_LEVEL", string(cfg.Severity))
	cfg.Severity = LogSeverity(sev)

	// Mode: console, file, kafka (comma-separated)
	modeEnv := getEnv("LOG_MODE", "console")
	modeList := strings.Split(modeEnv, ",")
	var modes []LogModeType
	for _, m := range modeList {
		switch strings.ToLower(strings.TrimSpace(m)) {
		case "console":
			modes = append(modes, ConsoleMode)
		case "file":
			modes = append(modes, FileMode)
			// future: case "kafka": modes = append(modes, KafkaMode)
		}
	}
	cfg.Mode = modes

	// File logging options
	cfg.FilePath = getEnv("LOG_FILE_PATH", cfg.FilePath)
	cfg.MaxSizeMB = getEnvAsInt("LOG_MAX_SIZE_MB", cfg.MaxSizeMB)
	cfg.MaxBackups = getEnvAsInt("LOG_MAX_BACKUPS", cfg.MaxBackups)
	cfg.MaxAgeDays = getEnvAsInt("LOG_MAX_AGE_DAYS", cfg.MaxAgeDays)

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
