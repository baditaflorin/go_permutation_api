package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

// Cache implements a simple in-memory cache with TTL
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewCache creates a new cache with the given TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go c.cleanup()

	return c
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in the cache
func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// cleanup removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// generateCacheKey creates a cache key from elements
func generateCacheKey(elements []string) string {
	// Sort elements for consistent key generation
	sorted := make([]string, len(elements))
	copy(sorted, elements)
	sort.Strings(sorted)

	// Create hash from sorted elements
	data, _ := json.Marshal(sorted)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
