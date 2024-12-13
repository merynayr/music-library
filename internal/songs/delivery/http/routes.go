package http

import (
	"music-library/internal/songs"

	"github.com/gin-gonic/gin"
)

func MapSongsRoutes(songsGroup *gin.RouterGroup, h songs.Handlers) {
	songsGroup.POST("/create", h.AddSong())
	songsGroup.DELETE("/delete/:id", h.DeleteSong())
	songsGroup.GET("", h.GetSongs())
	songsGroup.GET(":id/text", h.GetSongText())
	songsGroup.PUT("/:id", h.UpdateSong())
}
