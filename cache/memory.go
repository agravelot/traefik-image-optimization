package cache

import (
	"fmt"
	"sync"
	"time"
)

// MemoryCache in-memory cache system struct.
type MemoryCache struct {
	mtx sync.RWMutex
	m   map[string][]byte
}

// Get return cached image with given key.
func (c *MemoryCache) Get(key string) ([]byte, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	v, ok := c.m[key]
	if !ok {
		return nil, fmt.Errorf("no result found with key = %s", key)
	}

	return v, nil
}

// Set add a value into in-memory with custom expiry.
func (c *MemoryCache) Set(key string, v []byte, expiry time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[key] = v

	time.AfterFunc(expiry, func() {
		c.delete(key)
	})

	return nil
}

func (c *MemoryCache) delete(key string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.m, key)
}
