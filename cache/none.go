package cache

import (
	"fmt"
	"time"
)

// NoneCache dummy cache system.
type NoneCache struct{}

// Get always return nil with not found error.
func (c *NoneCache) Get(key string) ([]byte, error) {
	return nil, fmt.Errorf("no result found with key = %s", key)
}

// Set always return nil.
func (c *NoneCache) Set(_ string, _ []byte, _ time.Duration) error {
	return nil
}
