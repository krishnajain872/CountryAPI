package cache

import (
	"context"
	"time"

	"CountryAPI/pkg/types"
)

// Cache interface - The CONTRACT
type Cache interface {
	// Get retrieves a value by key
	// Returns: value and error (types.ErrCacheMiss if not found)
	Get(ctx context.Context, key string) (any, error)

	// Set stores a value with a time-to-live (TTL)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Delete removes a single entry
	Delete(ctx context.Context, key string) error

	// Clear removes all entries
	Clear(ctx context.Context) error

	// Stats returns cache statistics
	Stats() types.Stats
}
