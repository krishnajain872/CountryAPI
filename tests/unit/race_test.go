package unit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 1000)
	ctx := context.Background()
	
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100
	
	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Set(ctx, key, fmt.Sprintf("value-%d-%d", id, j), 0)
			}
		}(i)
	}
	
	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Get(ctx, key)
			}
		}(i)
	}
	
	wg.Wait()
	
	stats := cache.Stats()
	if stats.Size < 0 {
		t.Error("Cache size should not be negative")
	}
}

func TestMemoryCache_ConcurrentReadWrite(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 1000)
	ctx := context.Background()
	
	var wg sync.WaitGroup
	iterations := 1000
	
	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			cache.Set(ctx, "shared-key", i, 0)
		}
	}()
	
	// Reader goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				cache.Get(ctx, "shared-key")
			}
		}()
	}
	
	wg.Wait()
	
	// Should not crash or panic
	t.Log("Concurrent read/write test completed successfully")
}
// TestMemoryCache_ConcurrentSet tests concurrent writes
func TestMemoryCache_ConcurrentSet(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 10000)
	defer cache.Stop()
	ctx := context.Background()

	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Set(ctx, key, fmt.Sprintf("value-%d-%d", id, j), 0)
			}
		}(i)
	}

	wg.Wait()

	stats := cache.Stats()
	if stats.Size < 0 {
		t.Error("Cache size should not be negative")
	}
	
	t.Logf("Final cache size: %d", stats.Size)
}

// TestMemoryCache_ConcurrentGet tests concurrent reads
func TestMemoryCache_ConcurrentGet(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 10000)
	defer cache.Stop()
	ctx := context.Background()

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		cache.Set(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), 0)
	}

	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				cache.Get(ctx, fmt.Sprintf("key-%d", j%100))
			}
		}()
	}

	wg.Wait()
	
	stats := cache.Stats()
	if stats.Hits == 0 {
		t.Error("Should have recorded hits")
	}
	
	t.Logf("Hits: %d", stats.Hits)
}

// TestMemoryCache_ConcurrentDelete tests concurrent deletes
func TestMemoryCache_ConcurrentDelete(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 10000)
	defer cache.Stop()
	ctx := context.Background()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Set(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), 0)
	}

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d", id*10+j)
				cache.Delete(ctx, key)
			}
		}(i)
	}

	wg.Wait()
	t.Log("Concurrent delete test completed successfully")
}

// TestMemoryCache_ConcurrentStats tests stats under concurrent access
func TestMemoryCache_ConcurrentStats(t *testing.T) {
	cache := NewMemoryCache(5*time.Minute, 10000)
	defer cache.Stop()
	ctx := context.Background()

	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Set(ctx, fmt.Sprintf("key-%d-%d", id, j), "value", 0)
			}
		}(i)
	}

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Get(ctx, fmt.Sprintf("key-%d-%d", id, j))
			}
		}(i)
	}

	// Stats readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stats := cache.Stats()
				if stats.Size < 0 {
					t.Error("Invalid stats during concurrent access")
				}
			}
		}()
	}

	wg.Wait()
	t.Log("Concurrent stats test completed successfully")
}
