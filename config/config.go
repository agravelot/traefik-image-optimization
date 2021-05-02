package config

type ImaginaryProcessorConfig struct {
	Url string `json:"url" yaml:"url" toml:"url"`
}

type RedisCacheConfig struct {
	Url string `json:"url" yaml:"url" toml:"url"`
}

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
