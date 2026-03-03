package setup

import (
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"github.com/mahmud-off/auth/pkg/psql"
	"github.com/spf13/viper"
)

const (
	CONFIG_PATH = "configs"
	CONFIG_NAME = "main"
)

type Config struct {
	DB            psql.PostgresConnectionConfig
	HashSalt      string
	Port          string
	JSONFormatter bool
	Level         string
	ShowMethod    bool
}

func ParseConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("db", &cfg.DB)
	if err != nil {
		return nil, err
	}

	cfg.HashSalt = os.Getenv("HASH_SALT")

	viper.AddConfigPath(CONFIG_PATH)
	viper.SetConfigName(CONFIG_NAME)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg.Port = strconv.Itoa(viper.GetInt("server.port"))

	cfg.JSONFormatter = viper.GetBool("logger.JSONFormatter")
	cfg.Level = viper.GetString("logger.Level")
	cfg.ShowMethod = viper.GetBool("logger.ShowMethod")

	return &cfg, nil
}
