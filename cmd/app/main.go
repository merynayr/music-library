package main

import (
	"music-library/config"
	"music-library/internal/server"
	"music-library/pkg/db/postgres"
	"music-library/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// init logger
	log := logger.GetLogger()

	// init config
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	cfg := config.MustLoad()

	// init storage
	psqlDB, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatal("failed to init storage", err)
	}
	defer psqlDB.Close()

	log.Info("starting storage")

	// init server
	s := server.NewServer(cfg, psqlDB, log)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}

	log.Info("Server stopped")
}
