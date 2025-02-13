// cache package is a universal cache handler designed to store and retrieve
// various types of data using Redis. This package uses generics to support
// multiple types of data seamlessly, allowing for type-safe operations
// with caching.
//
// Example usage with a generic type:
//
//	var cache repository.Cache[domain.GetTransactionStatusResponse]
//	data, err := cache.Get(ctx, "key")
//	if err != nil {
//	}
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/redis/go-redis/v9"
)

// preKey is used to distinguish the key from other services
const preKey = "MAG:"

var ErrCacheDisable = errors.New("cache disabled")
var ErrNotFound = errors.New("data not found")

// CacheStorer is a generic interface for cache operations.
type CacheStorer[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
}

// Cache is a generic cache handler for any type of data.
type Cache[T any] struct {
	rds     *redis.Client
	cacheON bool
}

// NewCache creates a new instance of Cache.
func NewCache[T any](rds *redis.Client, active bool) CacheStorer[T] {
	return &Cache[T]{
		rds:     rds,
		cacheON: active,
	}
}

func (c *Cache[T]) Get(ctx context.Context, key string) (T, error) {
	ctx, span := observ.GetTracer().Start(ctx, "cache.Get")
	defer span.End()

	var result T
	if !c.cacheON {
		return result, ErrCacheDisable
	}

	key = preKey + key
	res, err := c.rds.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, ErrNotFound
		}
		return result, err
	}

	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *Cache[T]) Set(ctx context.Context, key string, value T, expiration time.Duration) error {
	ctx, span := observ.GetTracer().Start(ctx, "cache.Set")
	defer span.End()

	if !c.cacheON {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	key = preKey + key
	return c.rds.Set(ctx, key, string(data), expiration).Err()
}
