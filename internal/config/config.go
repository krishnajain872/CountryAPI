// internal/config/config.go
package config

import "time"

// ====================
// Config Structs
// ====================

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type CacheConfig struct {
	TTL     time.Duration
	MaxSize int
}

type ExternalAPIConfig struct {
	RestCountriesURL string
	Timeout          time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
}

// ====================
// Defaults
// ====================

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:         "localhost",
		Port:         "8000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

func defaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		TTL:     5 * time.Minute,
		MaxSize: 1000,
	}
}

func defaultExternalAPIConfig() *ExternalAPIConfig {
	return &ExternalAPIConfig{
		RestCountriesURL: "https://restcountries.com/v3.1",
		Timeout:          10 * time.Second,
	}
}

func defaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:  "info",
		Format: "json",
	}
}
