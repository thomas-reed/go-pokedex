package pokecache

import (
	"time"
	"sync"
)

const (
	reapingInterval = 5 * time.Minute
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	data map[string]cacheEntry
	mu sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	if interval <= 0 {
		interval = reapingInterval
	}
	c := &Cache{
		data: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Extend(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.data[key]
	if !ok {
		return
	}
	entry.createdAt = time.Now()
	c.data[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			age := time.Since(entry.createdAt)
			if age > c.interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}