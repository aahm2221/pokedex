package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	mu       *sync.Mutex
	caches   map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		caches:   make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
		interval: interval,
	}
	go c.reapLoop()
	return &c
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.caches {
			if time.Since(val.createdAt) > c.interval {
				delete(c.caches, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.caches[key] = cacheEntry{val: val, createdAt: time.Now()}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cache, exists := c.caches[key]
	if exists {
		return cache.val, true
	} else {
		return nil, false
	}
}
