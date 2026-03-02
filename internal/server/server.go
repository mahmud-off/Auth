package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahmud-off/auth/pkg/logger"
)

type Server struct {
	httpSever *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {

	s.httpSever = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpSever.ListenAndServe()
}

func (s *Server) GracefulShutdown(ctx context.Context) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("SERVER SHUTTING DOWN")

	err := s.httpSever.Shutdown(ctx)
	if err != nil {
		logger.Errorf("server shutting down error: %s", err.Error())
	} else {
		logger.Info("server gracefully stopped")
	}
}
