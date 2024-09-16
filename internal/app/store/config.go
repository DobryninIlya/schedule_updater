package store

type Config struct {
	DatabaseURL string `toml:"database_url"`
}

// NewConfig is deprecated
func NewConfig() *Config {
	return &Config{}
}
