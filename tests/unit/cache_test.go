// ========================================
// FILE: pkg/cache_data/memory_cache_test.go
// Unit tests for cache_data implementation
// ========================================
package unit

import (
	"context"
	"fmt"
	"testing"
	"time"
    cache_pkg "github.com/krishnajain872/country-search-api-cache/pkg/cache"
	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

// TestNewMemoryCache tests cache_data initialization
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
			cache_data := cache_pkg.NewMemoryCache(tt.ttl, tt.maxSize)
			
			if cache_data == nil {
				t.Fatal("cache_pkg.NewMemoryCache returned nil")
			}
			
			if cache_data.TTL() != tt.TTL() {
				t.Errorf("Expected ttl %v, got %v", tt.ttl, cache_data.ttl)
			}
			
			if cache_data.maxSize != tt.maxSize {
				t.Errorf("Expected maxSize %d, got %d", tt.maxSize, cache_data.maxSize)
			}
			
			if cache_data.data == nil {
				t.Error("cache_data data map not initialized")
			}
			
			cache_data.Stop()
		})
	}
}

// TestMemoryCache_Set tests the Set operation
func TestMemoryCache_Set(t *testing.T) {
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 100)
	defer cache_data.Stop()
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
			err := cache_data.Set(ctx, tt.key, tt.value, tt.ttl)
			
			if err != nil {
				t.Errorf("Set failed: %v", err)
			}
			
			stats := cache_data.Stats()
			if stats.Size == 0 {
				t.Error("cache_data size should be greater than 0 after Set")
			}
		})
	}
}

// TestMemoryCache_Get tests the Get operation
func TestMemoryCache_Get(t *testing.T) {
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 100)
	defer cache_data.Stop()
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
			cache_data.Set(ctx, tt.setupKey, tt.setupVal, 0)
			
			val, err := cache_data.Get(ctx, tt.getKey)
			
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
	cache_data := cache_pkg.NewMemoryCache(100*time.Millisecond, 100)
	defer cache_data.Stop()
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
			cache_data.Set(ctx, key, "value", tt.ttl)
			
			if tt.waitTime > 0 {
				time.Sleep(tt.waitTime)
			}
			
			_, err := cache_data.Get(ctx, key)
			
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
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 100)
	defer cache_data.Stop()
	ctx := context.Background()

	// Setup
	cache_data.Set(ctx, "key1", "value1", 0)
	
	// Verify exists
	_, err := cache_data.Get(ctx, "key1")
	if err != nil {
		t.Fatal("Setup failed: key should exist")
	}
	
	// Delete
	err = cache_data.Delete(ctx, "key1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	
	// Verify deleted
	_, err = cache_data.Get(ctx, "key1")
	if err != types.ErrCacheMiss {
		t.Error("Key should be deleted")
	}
	
	// Delete non-existing key (should not error)
	err = cache_data.Delete(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Delete non-existent key should not error: %v", err)
	}
}

// TestMemoryCache_Clear tests the Clear operation
func TestMemoryCache_Clear(t *testing.T) {
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 100)
	defer cache_data.Stop()
	ctx := context.Background()

	// Add multiple entries
	for i := 0; i < 10; i++ {
		cache_data.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 0)
	}
	
	stats := cache_data.Stats()
	if stats.Size != 10 {
		t.Errorf("Expected 10 entries, got %d", stats.Size)
	}
	
	// Clear
	err := cache_data.Clear(ctx)
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	
	// Verify cleared
	stats = cache_data.Stats()
	if stats.Size != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", stats.Size)
	}
}

// TestMemoryCache_Stats tests the Stats operation
func TestMemoryCache_Stats(t *testing.T) {
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 100)
	defer cache_data.Stop()
	ctx := context.Background()

	// Initial stats
	stats := cache_data.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Size != 0 {
		t.Error("Initial stats should be zero")
	}

	// Add entry and get (hit)
	cache_data.Set(ctx, "key1", "value1", 0)
	cache_data.Get(ctx, "key1")
	
	stats = cache_data.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Size != 1 {
		t.Errorf("Expected size 1, got %d", stats.Size)
	}

	// Get non-existent (miss)
	cache_data.Get(ctx, "key2")
	
	stats = cache_data.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
}

// TestMemoryCache_MaxSize tests size limit enforcement
func TestMemoryCache_MaxSize(t *testing.T) {
	maxSize := 10
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, maxSize)
	defer cache_data.Stop()
	ctx := context.Background()

	// Fill cache_data to max
	for i := 0; i < maxSize; i++ {
		cache_data.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 0)
	}
	
	stats := cache_data.Stats()
	if stats.Size != int64(maxSize) {
		t.Errorf("Expected size %d, got %d", maxSize, stats.Size)
	}

	// Add one more (should evict oldest)
	cache_data.Set(ctx, "new-key", "new-value", 0)
	
	stats = cache_data.Stats()
	if stats.Size > int64(maxSize) {
		t.Errorf("cache_data exceeded max size: %d > %d", stats.Size, maxSize)
	}
	
	if stats.Evictions < 1 {
		t.Error("Expected at least one eviction")
	}
}

// TestMemoryCache_LRU tests LRU eviction
func TestMemoryCache_LRU(t *testing.T) {
	cache_data := cache_pkg.NewMemoryCache(5*time.Minute, 3)
	defer cache_data.Stop()
	ctx := context.Background()

	// Add 3 entries
	cache_data.Set(ctx, "key1", "value1", 0)
	time.Sleep(10 * time.Millisecond)
	cache_data.Set(ctx, "key2", "value2", 0)
	time.Sleep(10 * time.Millisecond)
	cache_data.Set(ctx, "key3", "value3", 0)

	// Access key1 to make it recently used
	cache_data.Get(ctx, "key1")
	time.Sleep(10 * time.Millisecond)

	// Add new entry (should evict key2, the oldest unaccessed)
	cache_data.Set(ctx, "key4", "value4", 0)

	// Check that key1 and key3 still exist
	_, err1 := cache_data.Get(ctx, "key1")
	_, err3 := cache_data.Get(ctx, "key3")
	_, err2 := cache_data.Get(ctx, "key2")

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
	cache_data := cache_pkg.NewMemoryCache(100*time.Millisecond, 100)
	defer cache_data.Stop()
	ctx := context.Background()

	// Add entries with short TTL
	for i := 0; i < 5; i++ {
		cache_data.Set(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 100*time.Millisecond)
	}

	stats := cache_data.Stats()
	initialSize := stats.Size

	// Wait for cleanup (runs every minute, but entries expire in 100ms)
	time.Sleep(150 * time.Millisecond)
	
	// Manually trigger cleanup by accessing expired entry
	cache_data.Get(ctx, "key0")

	// Wait a bit more for cleanup goroutine
	time.Sleep(100 * time.Millisecond)

	// All entries should be expired
	for i := 0; i < 5; i++ {
		_, err := cache_data.Get(ctx, fmt.Sprintf("key%d", i))
		if err != types.ErrCacheMiss {
			t.Errorf("Entry key%d should be expired", i)
		}
	}

	stats = cache_data.Stats()
	t.Logf("Initial size: %d, Final size: %d, Evictions: %d", 
		initialSize, stats.Size, stats.Evictions)
}

 