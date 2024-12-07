package server

import (
	"context"
	"database/sql"
	"music-library/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	gin    *gin.Engine
	cfg    *config.Config
	db     *sql.DB
	logger *logrus.Logger
}

func NewServer(cfg *config.Config, db *sql.DB, logger *logrus.Logger) *Server {
	return &Server{gin: gin.New(),
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
		s.logger.Infof("Server is listening on PORT: http://%s", s.cfg.Server.Address)
		if err := server.ListenAndServe(); err != nil {
			s.logger.Fatalf("Error starting Server: %v", err)
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
