package main

import (
	"errors"
	"net/http"
	"os"

	"music-library/internal/config"
	httpserver "music-library/internal/http-server"
	"music-library/internal/logger"
	"music-library/internal/server"
	"music-library/internal/service"
	"music-library/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.GetLogger()

	// init config

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	cfg := config.MustLoad()

	// init storage

	storage, err := storage.InitDB(storage.Config(cfg.Storage))
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}
	log.Info("starting storage")

	service, err := service.New(cfg, storage)
	if err != nil {
		log.Error("failed to init service", err)
		os.Exit(1)
	}

	h := httpserver.New(cfg, service)

	srv := server.New(cfg, h)

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to init http server", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	<-quit

	log.Info("Server stopped")
}
