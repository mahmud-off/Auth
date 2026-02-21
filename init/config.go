package setup

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/mahmud-off/auth/pkg/psql"
)

type Config struct {
	DB       psql.PostgresConnectionConfig
	HashSalt string
}

func ParseConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("db", &cfg.DB)
	if err != nil {
		return nil, err
	}

	cfg.HashSalt = os.Getenv("HASH_SALT")

	return &cfg, nil
}
