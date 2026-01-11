// internal/config/storage.go
package config

import "sync"

// ====================
// Singleton storage
// ====================

var (
	instance *Config
	mu       sync.RWMutex
)

// Config is the central struct
type Config struct {
	Server   *ServerConfig
	Cache    *CacheConfig
	External *ExternalAPIConfig
	Logger   *LoggerConfig
}

// Store saves config in memory
func Store(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	instance = cfg
}

// Get retrieves config from memory
func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}
