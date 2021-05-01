package config

type ImaginaryConfig struct {
	Url string `json:"url" yaml:"url" toml:"url"`
}

// Config the plugin configuration.
type Config struct {
	Processor string          `json:"processor" yaml:"processor" toml:"processor"`
	Imaginary ImaginaryConfig `json:"imaginary,omitempty" yaml:"imaginary,omitempty" toml:"imaginary,omitempty"`
}
