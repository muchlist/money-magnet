package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/spend/port"
	"github.com/muchlist/moneymagnet/pkg/cache"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

type ETagCache struct {
	cache       cache.CacheStorer[int64]
	defDuration time.Duration
	log         mlogger.Logger
}

// make sure the implementation satisfies the interface
var _ port.ETagStorer = (*ETagCache)(nil)

func NewETagCache(cache cache.CacheStorer[int64], defDuration time.Duration, logger mlogger.Logger) *ETagCache {
	return &ETagCache{
		cache:       cache,
		defDuration: defDuration,
		log:         logger,
	}
}

func (r *ETagCache) GetTagByPocketID(ctx context.Context, pocketID string) (int64, error) {
	value, err := r.cache.Get(ctx, pocketID)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("key value database error: %w", err)
	}

	return value, err
}

func (r *ETagCache) SetTagByPocketID(ctx context.Context, pocketID string, value int64) error {
	err := r.cache.Set(ctx, pocketID, value, r.defDuration)
	if err != nil {
		return fmt.Errorf("key value database error: %w", err)
	}
	return nil
}
