// File: internal/cache/memory_cache_test.go
package cache

import (
	"context"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)
	defer c.Stop()

	err := c.Set(ctx, "key1", "value1", 0)
	assert.NoError(t, err)

	val, err := c.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestMemoryCache_GetNonExistentKey(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)
	defer c.Stop()

	val, err := c.Get(ctx, "missing")
	assert.Error(t, err)
	assert.Nil(t, val)
	assert.Equal(t, ErrCacheMiss, err)
}

func TestMemoryCache_TTLExpiration(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(50*time.Millisecond, 10)
	defer c.Stop()

	err := c.Set(ctx, "key1", "value1", 50*time.Millisecond)
	assert.NoError(t, err)

	time.Sleep(60 * time.Millisecond)

	val, err := c.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Nil(t, val)
	assert.Equal(t, ErrCacheMiss, err)
}

func TestMemoryCache_Delete(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)
	defer c.Stop()

	_ = c.Set(ctx, "key1", "value1", 0)
	err := c.Delete(ctx, "key1")
	assert.NoError(t, err)

	val, err := c.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Equal(t, ErrCacheMiss, err)
	assert.Nil(t, val)
}

func TestMemoryCache_Clear(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)
	defer c.Stop()

	_ = c.Set(ctx, "key1", "value1", 0)
	_ = c.Set(ctx, "key2", "value2", 0)

	err := c.Clear(ctx)
	assert.NoError(t, err)

	val, err := c.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Nil(t, val)
	val, err = c.Get(ctx, "key2")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestMemoryCache_Stats(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)
	defer c.Stop()

	// Initially stats should be zero
	stats := c.Stats()
	assert.Equal(t, int64(0), stats.Hits)
	assert.Equal(t, int64(0), stats.Misses)
	assert.Equal(t, int64(0), stats.Evictions)
	assert.Equal(t, int64(0), stats.Size)

	_ = c.Set(ctx, "key1", "value1", 0)
	_, _ = c.Get(ctx, "key1") // hit
	_, _ = c.Get(ctx, "missing") // miss

	stats = c.Stats()
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, int64(0), stats.Evictions)
	assert.Equal(t, int64(1), stats.Size)
}

func TestMemoryCache_EvictionWhenMaxSizeReached(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 2) // maxSize = 2
	defer c.Stop()

	_ = c.Set(ctx, "key1", "v1", 0)
	time.Sleep(1 * time.Millisecond) // ensure time difference for LRU
	_ = c.Set(ctx, "key2", "v2", 0)
	time.Sleep(1 * time.Millisecond)
	_ = c.Set(ctx, "key3", "v3", 0) // should evict "key1"

	val, err := c.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Equal(t, ErrCacheMiss, err)
	assert.Nil(t, val)

	val, err = c.Get(ctx, "key2")
	assert.NoError(t, err)
	assert.Equal(t, "v2", val)

	val, err = c.Get(ctx, "key3")
	assert.NoError(t, err)
	assert.Equal(t, "v3", val)
}

func TestMemoryCache_CleanupRemovesExpired(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(10*time.Millisecond, 10)
	defer c.Stop()

	_ = c.Set(ctx, "key1", "v1", 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond) // wait for cleanup to trigger

	// Wait for cleanup goroutine
	time.Sleep(100 * time.Millisecond)

	val, err := c.Get(ctx, "key1")
	assert.Error(t, err)
	assert.Equal(t, ErrCacheMiss, err)
	assert.Nil(t, val)

	stats := c.Stats()
	assert.Equal(t, int64(1), stats.Evictions)
}

func TestMemoryCache_Stop(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache(1*time.Minute, 10)

	_ = c.Set(ctx, "key1", "v1", 0)
	c.Stop()

	// Setting after stop should still work
	err := c.Set(ctx, "key2", "v2", 0)
	assert.NoError(t, err)

	val, err := c.Get(ctx, "key2")
	assert.NoError(t, err)
	assert.Equal(t, "v2", val)
}
