package store

type Config struct {
	DatabaseURL   string `toml:"database_url"`
	ServiceApiKey string `toml:"service_api_key"`
	DataGateway   string `toml:"data_gateway"`
	Debug         bool   `toml:"debug"`
}

// NewConfig is deprecated
func NewConfig() *Config {
	return &Config{}
}
