package pokecache

import (
	"sync"
	"time"
)

type CacheStruct struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) CacheStruct {
	var c CacheStruct
	c.cache = make(map[string]cacheEntry)
	c.mu = &sync.Mutex{}

	//go c.reapLoop(interval)

	return c
}

func (c *CacheStruct) Add(newKey string, newVal []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var newCacheEntry cacheEntry
	newCacheEntry.createdAt = time.Now()
	newCacheEntry.val = newVal
	c.cache[newKey] = newCacheEntry
}

func (c *CacheStruct) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.cache[key]
	return entry.val, found
}

func (c *CacheStruct) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *CacheStruct) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cache {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.cache, k)
		}
	}
}
