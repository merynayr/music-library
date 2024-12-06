package handlers

import (
	"music-library/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongHandler struct {
	music service.Music
}

func NewSongHandler(music service.Music) *SongHandler {
	return &SongHandler{music: music}
}

func (h *SongHandler) GetSongs(c *gin.Context) {
	c.JSON(http.StatusOK, h.music.GetAllSongs())
}
