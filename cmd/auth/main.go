package main

import (
	"context"

	_ "github.com/lib/pq"
	setup "github.com/mahmud-off/auth/init"
	"github.com/mahmud-off/auth/internal/repository"
	"github.com/mahmud-off/auth/internal/server"
	"github.com/mahmud-off/auth/internal/service"
	"github.com/mahmud-off/auth/internal/transport/rest"
	"github.com/mahmud-off/auth/pkg/hash"
	"github.com/mahmud-off/auth/pkg/logger"
	"github.com/mahmud-off/auth/pkg/psql"
)

func main() {

	cfg, err := setup.ParseConfig()
	if err != nil {
		logger.Fatalf("Config parsing problem: %s", err.Error())
		return
	}

	logger.Init(&logger.LoggerConfig{
		JSONFormatter: cfg.JSONFormatter,
		Level:         cfg.Level,
		ShowMethod:    cfg.ShowMethod,
	})

	db, err := psql.NewPostgresDB(&cfg.DB)
	if err != nil {
		logger.Fatalf("PostgreSQL connection problem: %s", err.Error())
		return
	}

	hasher := hash.NewSHA1Hasher(cfg.HashSalt)

	UserRepo := repository.NewUsersRepository(db)
	InfoRepo := repository.NewInfoRepository(db)
	TokenRepo := repository.NewTokens(db)

	userService := service.NewUsersService(UserRepo, hasher, TokenRepo, []byte("some secret signature"))
	InfoService := service.NewInfoService(InfoRepo)

	handler := rest.NewHandler(userService, InfoService)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg.Port, handler.InitRoutes()); err != nil {
			logger.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("SERVER STARTED")

	srv.GracefulShutdown(context.Background())

	if err := db.Close(); err != nil {
		logger.Errorf("error closing database: %s", err.Error())
	}
}
