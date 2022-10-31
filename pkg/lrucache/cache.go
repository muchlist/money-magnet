package lrucache

import (
	lru "github.com/hashicorp/golang-lru"
)

type LRU struct {
	cache *lru.Cache
}

func NewLRUCache() CacheStorer {
	cache, _ := lru.New(200)
	return &LRU{
		cache: cache,
	}
}

// Get implements CacheStorer
func (l *LRU) Get(key string) (Payload, bool) {
	val, ok := l.cache.Get(key)
	if !ok {
		return Payload{}, false
	}

	return val.(Payload), true
}

// Set implements CacheStorer
func (l *LRU) Set(key string, value Payload) {
	l.cache.Add(key, value)
}

// Clear implements CacheStorer
func (l *LRU) Clear(key string) {
	l.cache.Remove(key)
}
