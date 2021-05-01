package cache

import (
	"fmt"
	"sync"
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

func (c *MemoryCache) Set(key string, v []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.m[key] = v
	return nil
}
