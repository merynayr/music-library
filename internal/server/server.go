package server

import (
	"music-library/internal/config"
	"music-library/internal/http-server"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, router *httpserver.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Address,
			Handler:      router.NewRouter(),
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
