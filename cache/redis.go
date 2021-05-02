package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisCache struct {
	client *redis.Client
}

func (c *RedisCache) Get(key string) ([]byte, error) {

	v, err := c.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		return nil, fmt.Errorf("no result found with key = %s", key)
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func (c *RedisCache) Set(key string, v []byte) error {
	return c.client.Set(ctx, key, v, 0).Err()
}
