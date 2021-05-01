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
	if conf.Cache == "memory" {
		return &MemoryCache{}, nil
	}

	// if conf.Processor == "redis" {
	// 	p, err := NewRedis(conf)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return p, nil
	// }

	return nil, fmt.Errorf("unable to resolve given cache %s", conf.Cache)
}
