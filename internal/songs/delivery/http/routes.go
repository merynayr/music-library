package http

import (
	"music-library/internal/songs"

	"github.com/gin-gonic/gin"
)

func MapSongsRoutes(songsGroup *gin.RouterGroup, h songs.Handlers) {
	songsGroup.POST("/create", h.AddSong())
}
