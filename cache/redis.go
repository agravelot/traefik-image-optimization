package cache

import "time"

// var ctx = context.Background()

type RedisCache struct {
	// client *redis.Client
}

func (c *RedisCache) Get(key string) ([]byte, error) {

	// v, err := c.client.Get(ctx, key).Bytes()

	// if err == redis.Nil {
	// 	return nil, fmt.Errorf("no result found with key = %s", key)
	// } else if err != nil {
	// 	return nil, err
	// }

	// return v, nil

	return []byte("unsafe not supported by yaegi"), nil
}

func (c *RedisCache) Set(key string, v []byte, expiry time.Duration) error {
	// return c.client.Set(ctx, key, v, expiry).Err()
	return nil
}
