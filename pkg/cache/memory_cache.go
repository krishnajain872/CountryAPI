// File: internal/cache/memory_cache.go
package cache

import (
	"context"
	"time"
	"sync"
	"errors"
)

var ErrCacheMiss = errors.New("cache: key not found")

type Stats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	Size      int64
}

type cacheEntry struct {
	value      any
	expiry     time.Time
	lastAccess time.Time
}

type MemoryCache struct {
	mu      sync.RWMutex
	data    map[string]*cacheEntry
	stats   Stats
	ttl     time.Duration
	maxSize int
	stopCh  chan struct{}
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(ttl time.Duration, maxSize int) *MemoryCache {
	c := &MemoryCache{
		data:    make(map[string]*cacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
		stopCh:  make(chan struct{}),
	}

	go c.cleanup()
	return c
}

// Get retrieves a value by key
func (c *MemoryCache) Get(ctx context.Context, key string) (any, error) {
	c.mu.RLock()
	entry, exists := c.data[key]
	c.mu.RUnlock()

	if !exists || time.Now().After(entry.expiry) {
		c.mu.Lock()
		c.stats.Misses++
		c.mu.Unlock()
		return nil, ErrCacheMiss
	}

	// Update last access safely
	c.mu.Lock()
	entry.lastAccess = time.Now()
	c.stats.Hits++
	c.mu.Unlock()

	return entry.value, nil
}

// Set stores a value with optional TTL
func (c *MemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.ttl
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict oldest if exceeding max size
	if len(c.data) >= c.maxSize && c.data[key] == nil {
		c.evictOldest()
	}

	c.data[key] = &cacheEntry{
		value:      value,
		expiry:     time.Now().Add(ttl),
		lastAccess: time.Now(),
	}

	c.stats.Size = int64(len(c.data))
	return nil
}

// Delete removes a single entry
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	c.stats.Size = int64(len(c.data))
	return nil
}

// Clear removes all entries
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*cacheEntry)
	c.stats.Size = 0
	return nil
}

// Stats returns cache statistics
func (c *MemoryCache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// cleanup runs periodically to remove expired entries
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCh:
			return
		}
	}
}

// removeExpired deletes entries that have expired
func (c *MemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.expiry) {
			delete(c.data, key)
			c.stats.Evictions++
		}
	}
	c.stats.Size = int64(len(c.data))
}

// evictOldest removes the least recently accessed entry (LRU)
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.data {
		if oldestKey == "" || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
		c.stats.Evictions++
	}
}

// Stop stops the cleanup goroutine
func (c *MemoryCache) Stop() {
	close(c.stopCh)
}
