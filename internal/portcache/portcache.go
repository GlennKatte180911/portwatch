// Package portcache provides a short-lived in-memory cache for port scan
// results, reducing redundant scans within a configurable TTL window.
package portcache

import (
	"sync"
	"time"
)

// Entry holds a cached scan result and its expiry time.
type Entry struct {
	Ports     []int
	ScannedAt time.Time
	expiresAt time.Time
}

// Cache stores the most recent scan result per host key.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL. A zero TTL disables caching (every
// lookup is treated as a miss).
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Set stores ports for key, overwriting any existing entry.
func (c *Cache) Set(key string, ports []int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	copy := make([]int, len(ports))
	_ = copy[:copy(copy, ports)]
	c.entries[key] = &Entry{
		Ports:     copy,
		ScannedAt: now,
		expiresAt: now.Add(c.ttl),
	}
}

// Get returns the cached entry for key and whether it is still valid.
func (c *Cache) Get(key string) (*Entry, bool) {
	if c.ttl <= 0 {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e, true
}

// Invalidate removes the cached entry for key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Purge removes all expired entries from the cache.
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, e := range c.entries {
		if now.After(e.expiresAt) {
			delete(c.entries, k)
		}
	}
}

// Len returns the number of entries currently held (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
