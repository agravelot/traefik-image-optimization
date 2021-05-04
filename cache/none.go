package cache

import (
	"fmt"
	"time"
)

type NoneCache struct {
}

func (c *NoneCache) Get(key string) ([]byte, error) {
	return nil, fmt.Errorf("no result found with key = %s", key)
}

func (c *NoneCache) Set(_ string, _ []byte, _ time.Duration) error {
	return nil
}
