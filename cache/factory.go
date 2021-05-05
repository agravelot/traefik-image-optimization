// Package cache provide caching systems for images.
package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/agravelot/imageopti/config"
)

// Cache Define cache system interface.
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte, expiry time.Duration) error
}

const defaultCacheExpiry = 100 * time.Second

// New is the cache factory to instantiate a new instance of cache.
func New(conf config.Config) (Cache, error) {
	// if conf.Processor == "redis" {
	// 	opt, err := redis.ParseURL(conf.Redis.URL)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	client := redis.NewClient(opt)

	// 	return &RedisCache{
	// 		client: client,
	// 	}, nil
	// }

	if conf.Cache == "file" {
		return newFileCache(conf.File.Path, defaultCacheExpiry)
	}

	if conf.Cache == "memory" {
		return &MemoryCache{
			m:   map[string][]byte{},
			mtx: sync.RWMutex{},
		}, nil
	}

	if conf.Cache == "none" || conf.Cache == "" {
		return &NoneCache{}, nil
	}

	return nil, fmt.Errorf("unable to resolve given cache %s", conf.Cache)
}
