package cache

import (
	"fmt"

	"github.com/agravelot/image_optimizer/config"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, v []byte) error
}

// TODO Use singleton patern ?
func New(conf config.Config) (Cache, error) {
	// if conf.Processor == "redis" {
	// 	opt, err := redis.ParseURL(conf.Redis.Url)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	client := redis.NewClient(opt)

	// 	return &RedisCache{
	// 		client: client,
	// 	}, nil
	// }

	if conf.Cache == "memory" {
		return &MemoryCache{}, nil
	}

	if conf.Cache == "none" {
		return &NoneCache{}, nil
	}

	return nil, fmt.Errorf("unable to resolve given cache %s", conf.Cache)
}
