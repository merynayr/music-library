package httpserver

import (
	"music-library/internal/config"
	"music-library/internal/http-server/handlers"
	"music-library/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg   *config.Config
	music service.Music
}

func New(cfg *config.Config, music service.Music) *Handler {
	return &Handler{
		cfg:   cfg,
		music: music,
	}
}

func (h *Handler) NewRouter() http.Handler {
	router := gin.Default()

	songHandler := handlers.NewSongHandler(h.music)
	router.GET("/songs", songHandler.GetSongs)

	return router
}
