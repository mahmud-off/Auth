package main

import (
	_ "github.com/lib/pq"
	setup "github.com/mahmud-off/auth/init"
	"github.com/mahmud-off/auth/internal/repository"
	"github.com/mahmud-off/auth/internal/server"
	"github.com/mahmud-off/auth/internal/service"
	"github.com/mahmud-off/auth/internal/transport/rest"
	"github.com/mahmud-off/auth/pkg/hash"
	"github.com/mahmud-off/auth/pkg/psql"
	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := setup.ParseConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := psql.NewPostgresDB(&cfg.DB)
	if err != nil {
		logrus.Fatal(err)
	}

	hasher := hash.NewSHA1Hasher(cfg.HashSalt)

	UserRepo := repository.NewUsersRepository(db)
	InfoRepo := repository.NewInfoRepository(db)
	TokenRepo := repository.NewTokens(db)

	userService := service.NewUsersService(UserRepo, hasher, TokenRepo, []byte("some secret signature"))
	InfoService := service.NewInfoService(InfoRepo)

	handler := rest.NewHandler(userService, InfoService)

	srv := new(server.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}
