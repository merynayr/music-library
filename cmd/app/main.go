package main

import (
	"log/slog"
	"os"

	"music-library/config"
	"music-library/internal/server"
	"music-library/pkg/db/postgres"
	"music-library/pkg/logger"
	"music-library/pkg/logger/sl"

	"github.com/joho/godotenv"
)

func main() {
	// init config
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("failed to init logger", sl.Err(err))
		os.Exit(1)
	}

	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)

	log.Info(
		"starting Project",
		"log level", cfg.Env,
	)

	// init storage
	psqlDB, err := postgres.InitDB(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer psqlDB.Close()

	log.Info("starting storage")

	// init server
	s := server.NewServer(cfg, psqlDB, log)
	if err = s.Run(); err != nil {
		log.Error("failed to start server")
	}

	log.Info("Server stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	slog.SetDefault(log)
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
