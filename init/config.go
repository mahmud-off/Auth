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
	SQLDB         psql.PostgresConnectionConfig
	HashSalt      string
	SQLPort       string
	JSONFormatter bool
	Level         string
	ShowMethod    bool
	RedisPort     string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func ParseConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("db", &cfg.SQLDB)
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

	cfg.SQLPort = strconv.Itoa(viper.GetInt("server.port"))

	cfg.JSONFormatter = viper.GetBool("logger.JSONFormatter")
	cfg.Level = viper.GetString("logger.Level")
	cfg.ShowMethod = viper.GetBool("logger.ShowMethod")

	cfg.RedisAddr = viper.GetString("redis.addr")
	cfg.RedisPort = strconv.Itoa(viper.GetInt("redis.port"))
	cfg.RedisPassword = viper.GetString("redis.password")
	cfg.RedisDB = viper.GetInt("redis.db")

	return &cfg, nil
}
