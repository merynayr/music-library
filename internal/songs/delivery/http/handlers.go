package http

import (
	"music-library/config"
	"music-library/internal/models"
	"music-library/internal/songs"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Songs handlers
type songsHandlers struct {
	cfg     *config.Config
	songsUC songs.UseCase
	logger  *logrus.Logger
}

// NewSongsHandlers Songs handlers constructor
func NewSongsHandlers(cfg *config.Config, songsUC songs.UseCase, logger *logrus.Logger) songs.Handlers {
	return &songsHandlers{cfg: cfg, songsUC: songsUC, logger: logger}
}

// @Router /songs/create [post]
func (h songsHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.Create"

		n := &models.Song{}
		if err := c.Bind(n); err != nil {
			h.logger.Error("Failed to get users", op, err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		createdSongs, err := h.songsUC.Create(n)
		if err != nil {
			h.logger.Error("Failed to get users", op, err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusCreated, createdSongs)
	}
}

// @Router /songs/hello [get]
func (h songsHandlers) Hello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, "Hello my friend")
	}
}
