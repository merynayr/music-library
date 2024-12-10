package http

import (
	"encoding/json"
	"fmt"
	"io"
	"music-library/config"
	"music-library/internal/models"
	"music-library/internal/songs"
	resp "music-library/pkg/api"
	"music-library/pkg/logger/sl"
	"strconv"

	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Songs handlers
type songsHandlers struct {
	cfg     *config.Config
	songsUC songs.UseCase
	log     *slog.Logger
}

// NewSongsHandlers Songs handlers constructor
func NewSongsHandlers(cfg *config.Config, songsUC songs.UseCase, logger *slog.Logger) songs.Handlers {
	return &songsHandlers{cfg: cfg, songsUC: songsUC, log: logger}
}

type Request struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type Response struct {
	resp.Response
}

// @Router /songs/create [post]
func (h songsHandlers) AddSong() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.AddSong"

		log := h.log.With(
			slog.String("op", op),
		)

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed to bind JSON: ", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		log.Debug("Request body decoded", slog.Any("request", req))

		songDetail, err := GetSongInfoFromAPI(h, req)
		if err != nil {
			log.Error("Failed to get song info from external API", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid info request body"})
			return
		}
		log.Debug("Retrieved song details from external API", slog.Any("request", songDetail))
		songDetail.SongName = req.Song

		group := &models.Group{Name: req.Group}

		createdSongs, err := h.songsUC.AddGroupWithSongsTx(group, &songDetail)
		if err != nil {
			log.Error("Failed to add group and songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add group and songs"})
			return
		}

		c.JSON(http.StatusCreated, createdSongs)
	}
}

func GetSongInfoFromAPI(h songsHandlers, req Request) (models.Song, error) {
	op := "songsHandlers.GetSongInfoFromAPI"

	log := h.log.With(
		slog.String("op", op),
	)
	group, song := req.Group, req.Song

	url := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", url.QueryEscape(group), url.QueryEscape(song))

	resp, err := http.Get(url)
	if err != nil {
		log.Error("Get API request", sl.Err(err))
		return models.Song{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Reading response body", sl.Err(err))
		return models.Song{}, err
	}

	var songDetail models.Song
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		log.Error("Unmarshalling JSON response", sl.Err(err))
		return models.Song{}, err
	}

	return songDetail, nil
}

// @Router /songs/delete/:id [delete]
func (h songsHandlers) DeleteSong() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.DeleteSong"

		log := h.log.With(
			slog.String("op", op),
		)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Error("Failed to parse id", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
			return
		}
		log.Debug("Request to delete song with ID", slog.Any("request", id))

		err = h.songsUC.DeleteSong(uint(id))
		if err != nil {
			log.Error("Failed to delete songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Song deleted"})
	}
}

// @Router /songs [get]
func (h songsHandlers) GetSongs() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.GetSongs"

		log := h.log.With(
			slog.String("op", op),
		)
		songs, err := h.songsUC.GetSongs()
		if err != nil {
			log.Error("Failed to delete songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
			return
		}

		c.JSON(http.StatusOK, songs)
	}
}
