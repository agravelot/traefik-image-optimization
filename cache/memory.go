package cache

import (
	"fmt"
	"sync"
	"time"
)

type MemoryCache struct {
	mtx sync.RWMutex
	m   map[string][]byte
}

func (c *MemoryCache) Get(key string) ([]byte, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	v, ok := c.m[key]
	if !ok {
		return nil, fmt.Errorf("no result found with key = %s", key)
	}
	return v, nil
}

// TODO Implement exipire logic
func (c *MemoryCache) Set(key string, v []byte, _ time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[key] = v
	return nil
}
