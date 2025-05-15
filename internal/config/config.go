package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort  string `envconfig:"APP_PORT"   default:":8080"`
	LogLevel string `envconfig:"LOG_LEVEL"  default:"info"` // debug | info | warn | error
	DBDSN    string `envconfig:"DB_DSN"     default:"postgres://user:pass@localhost:5432/people?sslmode=disable"`
	CacheTTL string `envconfig:"CACHE_TTL"  default:"24h"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
