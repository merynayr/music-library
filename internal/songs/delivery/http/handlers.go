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
	"music-library/pkg/utils"
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

// AddSong добавляет новую песню в библиотеку
// @Summary Add a new song to the library
// @Description Adds a new song with details like group, song name, release date, text, and link
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song details"
// @Success 201 {object} models.Song "Song created successfully"
// @Failure 400 {object} api.Response "Invalid input data"
// @Failure 500 {object} api.Response "Failed to call external API or insert data into database"
// @Router /api/songs/create [post]
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

// DeleteSong удаляет песню из базы данных
// @Summary Delete a song by ID
// @Description Deletes a song from the database by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 200 {object} string "Song deleted successfully"
// @Failure 404 {object} api.Response "Failed to delete song"
// @Failure 500 {object} api.Response "Failed to delete song"
// @Router /api/songs/delete/:id [delete]
func (h songsHandlers) DeleteSong() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.DeleteSong"

		log := h.log.With(
			slog.String("op", op),
		)

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			log.Error("Failed to parse id", sl.Err(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete song"})
			return
		}
		log.Debug("Request to delete song with ID", slog.Any("request", id))

		err = h.songsUC.DeleteSong(uint(id))
		if err != nil {
			log.Error("Failed to delete songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
	}
}

// GetSongs возвращает список песен с фильтрацией и пагинацией
// @Summary Get songs list with filtering and pagination
// @Description Retrieves a paginated list of songs with optional filtering based on group, song name, release date, text, and link
// @Tags songs
// @Param groupName query string false "Group name for filtering"
// @Param song query string false "Song name for filtering"
// @Param releaseDate query string false "Release date for filtering"
// @Param text query string false "Text for filtering"
// @Param link query string false "Link for filtering"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of songs per page" default(10)
// @Success 200 {object}} models.Song
// @Failure 404 {object} api.Response "Failed to get song"
// @Failure 500 {object} api.Response "Failed to get songs"
// @Router /api/songs [get]
func (h songsHandlers) GetSongs() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.GetSongs"

		log := h.log.With(
			slog.String("op", op),
		)

		pq, err := utils.GetPaginationFromCtx(c.Copy())
		if err != nil {
			log.Error("Failed to GetPaginationFromCtx", sl.Err(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get songs"})
			return
		}

		fq, err := utils.GetFilterFromCtx(c.Copy())
		if err != nil {
			log.Error("Failed to GetFilterFromCtx", sl.Err(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get songs"})
			return
		}
		songs, err := h.songsUC.GetSongs(pq, fq)
		if err != nil {
			log.Error("Failed to delete songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
			return
		}

		c.JSON(http.StatusOK, songs)
	}
}

// GetSongText возвращает текст песни с пагинацией по куплетам
// @Summary Get song text by verses with pagination
// @Description Retrieves the song's text, paginated by verses, based on the song's ID
// @Tags songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of verses per page" default(2)
// @Success 200 {object} string
// @Failure 404 {object} api.Response "Failed to get song`s text"
// @Failure 500 {object} api.Response "Failed to get song`s text"
// @Router /api/songs/{id}/text [get]
func (h songsHandlers) GetSongText() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.GetSongText"
		log := h.log.With(
			slog.String("op", op),
		)

		id := c.Param("id")

		pq, err := utils.GetPaginationFromCtx(c.Copy())
		if err != nil {
			log.Error("Failed to get song`s text", sl.Err(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get song`s text"})
			return
		}

		text, err := h.songsUC.GetSongText(id, pq)
		if err != nil {
			log.Error("Failed to get song`s text", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song`s text"})
			return
		}
		c.JSON(http.StatusOK, text)

	}
}

// UpdateSong обновляет данные о песне по её ID
// @Summary Update song details
// @Description Updates the song information by its ID. Only provided fields will be updated.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} string "Song updated successfully"
// @Failure 400 {object} api.Response "Invalid JSON data"
// @Failure 500 {object} api.Response "Failed to update song"
// @Router /api/songs/:id [put]
func (h songsHandlers) UpdateSong() gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "songsHandlers.UpdateSong"

		log := h.log.With(
			slog.String("op", op),
		)

		id := c.Param("id")

		var Data map[string]interface{}
		if err := c.ShouldBindJSON(&Data); err != nil {
			log.Error("Failed to bind JSON: ", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		err := h.songsUC.UpdateSong(id, Data)
		if err != nil {
			log.Error("Failed to update songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
	}
}
