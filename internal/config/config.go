// internal/config/config.go
package config

import "time"

// ====================
// Logger helper types
// ====================

// LogEnvType represents environment mode
type LogEnvType string

const (
	Development LogEnvType = "development"
	Staging     LogEnvType = "staging"
	Production  LogEnvType = "production"
)

// LogModeType represents logging output type
type LogModeType string

const (
	ConsoleMode LogModeType = "console"
	FileMode    LogModeType = "file"
	// KafkaMode LogModeType = "kafka" // future extension
)

// LogSeverity represents minimum severity of logs
type LogSeverity string

const (
	DebugSeverity    LogSeverity = "debug"
	InfoSeverity     LogSeverity = "info"
	WarnSeverity     LogSeverity = "warn"
	ErrorSeverity    LogSeverity = "error"
	CriticalSeverity LogSeverity = "critical"
)

// ====================
// Logger Config
// ====================

type LoggerConfig struct {
	Environment LogEnvType
	Mode        []LogModeType
	Severity    LogSeverity
	FilePath    string
	MaxSizeMB   int
	MaxBackups  int
	MaxAgeDays  int
}

// ====================
// Other Configs
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

// ====================
// Defaults
// ====================

func defaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Environment: Development,
		Mode:        []LogModeType{ConsoleMode},
		Severity:    InfoSeverity,
		FilePath:    "logs/app.log",
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  7,
	}
}

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
