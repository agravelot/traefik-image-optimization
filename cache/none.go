package cache

import (
	"fmt"
)

type NoneCache struct {
}

func (c *NoneCache) Get(key string) ([]byte, error) {
	return nil, fmt.Errorf("no result found with key = %s", key)
}

func (c *NoneCache) Set(key string, v []byte) error {
	return nil
}
