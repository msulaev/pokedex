package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	items    map[string]cacheEntry
	interval time.Duration
	stopCh   chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		items:    make(map[string]cacheEntry),
		interval: interval,
		stopCh:   make(chan struct{}),
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.createdAt) > c.interval {
		delete(c.items, key)
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.reap()
		}
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, entry := range c.items {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.items, key)
		}
	}
}

func (c *Cache) Stop() {
	close(c.stopCh)
}
