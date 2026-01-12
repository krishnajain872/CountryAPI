// ========================================
// FILE: tests/unit/cache_test.go
// Unit tests for cache implementation
// ========================================
package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/krishnajain872/country-search-api-cache/pkg/cache"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

// TestNewMemoryCache tests cache initialization
func TestNewMemoryCache(t *testing.T) {
	tests := []struct {
		name    string
		ttl     time.Duration
		maxSize int
	}{
		{"default config", 5 * time.Minute, 1000},
		{"short ttl", 1 * time.Minute, 100},
		{"large size", 10 * time.Minute, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cache.NewMemoryCache(tt.ttl, tt.maxSize)

			if c == nil {
				t.Fatal("NewMemoryCache returned nil")
			}

			// Test with public methods instead of accessing private fields
			stats := c.Stats()
			if stats.Size != 0 {
				t.Errorf("Initial cache size should be 0, got %d", stats.Size)
			}

			c.Stop()
		})
	}
}

// TestMemoryCache_Set tests the Set operation
func TestMemoryCache_Set(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 100)
	defer c.Stop()
	ctx := context.Background()

	tests := []struct {
		name  string
		key   string
		value interface{}
		ttl   time.Duration
	}{
		{"string value", "key1", "value1", 0},
		{"int value", "key2", 42, 0},
		{"struct value", "key3", types.Country{Name: "India"}, 0},
		{"custom ttl", "key4", "value4", 1 * time.Minute},
		{"nil value", "key5", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Set(ctx, tt.key, tt.value, tt.ttl)

			if err != nil {
				t.Errorf("Set failed: %v", err)
			}

			stats := c.Stats()
			if stats.Size == 0 {
				t.Error("cache size should be greater than 0 after Set")
			}
		})
	}
}

// TestMemoryCache_Get tests the Get operation
func TestMemoryCache_Get(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 100)
	defer c.Stop()
	ctx := context.Background()

	tests := []struct {
		name      string
		setupKey  string
		setupVal  interface{}
		getKey    string
		wantError bool
		wantValue interface{}
	}{
		{
			name:      "existing key",
			setupKey:  "exists",
			setupVal:  "value",
			getKey:    "exists",
			wantError: false,
			wantValue: "value",
		},
		{
			name:      "non-existing key",
			setupKey:  "exists",
			setupVal:  "value",
			getKey:    "notexists",
			wantError: true,
		},
		{
			name:      "nil value",
			setupKey:  "nil",
			setupVal:  nil,
			getKey:    "nil",
			wantError: false,
			wantValue: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Set(ctx, tt.setupKey, tt.setupVal, 0)

			val, err := c.Get(ctx, tt.getKey)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if err != types.ErrCacheMiss {
					t.Errorf("Expected ErrCacheMiss, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if val != tt.wantValue {
					t.Errorf("Expected %v, got %v", tt.wantValue, val)
				}
			}
		})
	}
}

// TestMemoryCache_TTL tests time-to-live functionality
func TestMemoryCache_TTL(t *testing.T) {
	c := cache.NewMemoryCache(100*time.Millisecond, 100)
	defer c.Stop()
	ctx := context.Background()

	tests := []struct {
		name       string
		ttl        time.Duration
		waitTime   time.Duration
		shouldFind bool
	}{
		{"immediate access", 100 * time.Millisecond, 0, true},
		{"before expiry", 200 * time.Millisecond, 100 * time.Millisecond, true},
		{"after expiry", 100 * time.Millisecond, 150 * time.Millisecond, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := fmt.Sprintf("ttl-%s", tt.name)
			c.Set(ctx, key, "value", tt.ttl)

			if tt.waitTime > 0 {
				time.Sleep(tt.waitTime)
			}

			_, err := c.Get(ctx, key)

			if tt.shouldFind && err != nil {
				t.Errorf("Expected to find key, got error: %v", err)
			}

			if !tt.shouldFind && err != types.ErrCacheMiss {
				t.Error("Expected key to be expired")
			}
		})
	}
}

// TestMemoryCache_Delete tests the Delete operation
func TestMemoryCache_Delete(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 100)
	defer c.Stop()
	ctx := context.Background()

	// Setup
	c.Set(ctx, "key1", "value1", 0)

	// Verify exists
	_, err := c.Get(ctx, "key1")
	if err != nil {
		t.Fatal("Setup failed: key should exist")
	}

	// Delete
	err = c.Delete(ctx, "key1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = c.Get(ctx, "key1")
	if err != types.ErrCacheMiss {
		t.Error("Key should be deleted")
	}

	// Delete non-existing key (should not error)
	err = c.Delete(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Delete non-existent key should not error: %v", err)
	}
}

// TestMemoryCache_Clear tests the Clear operation
func TestMemoryCache_Clear(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 100)
	defer c.Stop()
	ctx := context.Background()

	// Add multiple entries
	for i := 0; i < 10; i++ {
		c.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 0)
	}

	stats := c.Stats()
	if stats.Size != 10 {
		t.Errorf("Expected 10 entries, got %d", stats.Size)
	}

	// Clear
	err := c.Clear(ctx)
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	// Verify cleared
	stats = c.Stats()
	if stats.Size != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", stats.Size)
	}
}

// TestMemoryCache_Stats tests the Stats operation
func TestMemoryCache_Stats(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 100)
	defer c.Stop()
	ctx := context.Background()

	// Initial stats
	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Size != 0 {
		t.Error("Initial stats should be zero")
	}

	// Add entry and get (hit)
	c.Set(ctx, "key1", "value1", 0)
	c.Get(ctx, "key1")

	stats = c.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Size != 1 {
		t.Errorf("Expected size 1, got %d", stats.Size)
	}

	// Get non-existent (miss)
	c.Get(ctx, "key2")

	stats = c.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
}

// TestMemoryCache_MaxSize tests size limit enforcement
func TestMemoryCache_MaxSize(t *testing.T) {
	maxSize := 10
	c := cache.NewMemoryCache(5*time.Minute, maxSize)
	defer c.Stop()
	ctx := context.Background()

	// Fill cache to max
	for i := 0; i < maxSize; i++ {
		c.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 0)
	}

	stats := c.Stats()
	if stats.Size != int64(maxSize) {
		t.Errorf("Expected size %d, got %d", maxSize, stats.Size)
	}

	// Add one more (should evict oldest)
	c.Set(ctx, "new-key", "new-value", 0)

	stats = c.Stats()
	if stats.Size > int64(maxSize) {
		t.Errorf("cache exceeded max size: %d > %d", stats.Size, maxSize)
	}

	if stats.Evictions < 1 {
		t.Error("Expected at least one eviction")
	}
}

// TestMemoryCache_LRU tests LRU eviction
func TestMemoryCache_LRU(t *testing.T) {
	c := cache.NewMemoryCache(5*time.Minute, 3)
	defer c.Stop()
	ctx := context.Background()

	// Add 3 entries
	c.Set(ctx, "key1", "value1", 0)
	time.Sleep(10 * time.Millisecond)
	c.Set(ctx, "key2", "value2", 0)
	time.Sleep(10 * time.Millisecond)
	c.Set(ctx, "key3", "value3", 0)

	// Access key1 to make it recently used
	c.Get(ctx, "key1")
	time.Sleep(10 * time.Millisecond)

	// Add new entry (should evict key2, the oldest unaccessed)
	c.Set(ctx, "key4", "value4", 0)

	// Check that key1 and key3 still exist
	_, err1 := c.Get(ctx, "key1")
	_, err3 := c.Get(ctx, "key3")
	_, err2 := c.Get(ctx, "key2")

	if err1 != nil {
		t.Error("key1 should still exist (was recently accessed)")
	}
	if err3 != nil {
		t.Error("key3 should still exist")
	}
	if err2 != types.ErrCacheMiss {
		t.Error("key2 should have been evicted")
	}
}

// TestMemoryCache_Cleanup tests automatic cleanup
func TestMemoryCache_Cleanup(t *testing.T) {
	c := cache.NewMemoryCache(100*time.Millisecond, 100)
	defer c.Stop()
	ctx := context.Background()

	// Add entries with short TTL
	for i := 0; i < 5; i++ {
		c.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 100*time.Millisecond)
	}

	stats := c.Stats()
	initialSize := stats.Size

	// Wait for cleanup (runs every minute, but entries expire in 100ms)
	time.Sleep(150 * time.Millisecond)

	// Manually trigger cleanup by accessing expired entry
	c.Get(ctx, "key0")

	// Wait a bit more for cleanup goroutine
	time.Sleep(100 * time.Millisecond)

	// All entries should be expired
	for i := 0; i < 5; i++ {
		_, err := c.Get(ctx, fmt.Sprintf("key%d", i))
		if err != types.ErrCacheMiss {
			t.Errorf("Entry key%d should be expired", i)
		}
	}

	stats = c.Stats()
	t.Logf("Initial size: %d, Final size: %d, Evictions: %d",
		initialSize, stats.Size, stats.Evictions)
}
