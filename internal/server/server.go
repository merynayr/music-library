package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music-library/config"
	"music-library/pkg/logger/sl"

	"github.com/gin-gonic/gin"
)

type Server struct {
	gin    *gin.Engine
	cfg    *config.Config
	db     *sql.DB
	logger *slog.Logger
}

func NewServer(cfg *config.Config, db *sql.DB, logger *slog.Logger) *Server {
	return &Server{
		gin:    gin.New(),
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         s.cfg.Address,
		Handler:      s.gin,
		ReadTimeout:  s.cfg.Server.Timeout,
		WriteTimeout: s.cfg.Server.Timeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	go func() {
		s.logger.Info("Server is listening", "address", s.cfg.Server.Address)
		if err := server.ListenAndServe(); err != nil {
			s.logger.Error("Error starting Server", sl.Err(err))
		}
	}()

	if err := s.MapHandlers(s.gin); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return server.Shutdown(ctx)
}
