// Package config provide configurations structs for imageopti middleware.
package config

// ImaginaryProcessorConfig define imaginary image processor configurations.
type ImaginaryProcessorConfig struct {
	URL string `json:"url" yaml:"url" toml:"url"`
}

// RedisCacheConfig define redis cache system configurations.
type RedisCacheConfig struct {
	URL string `json:"url" yaml:"url" toml:"url"`
}

// FileCacheConfig define file cache system configurations.
type FileCacheConfig struct {
	Path string `json:"path" yaml:"path" toml:"path"`
}

// Config the plugin configuration.
type Config struct {
	Processor string                   `json:"processor" yaml:"processor" toml:"processor"`
	Imaginary ImaginaryProcessorConfig `json:"imaginary,omitempty" yaml:"imaginary,omitempty" toml:"imaginary,omitempty"`
	// Cache
	Cache string           `json:"cache" yaml:"cache" toml:"cache"`
	Redis RedisCacheConfig `json:"redis,omitempty" yaml:"redis,omitempty" toml:"redis,omitempty"`
	File  FileCacheConfig  `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`
}
