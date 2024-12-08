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
	"time"

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
			c.JSON(http.StatusBadRequest, resp.Error(err))
			return
		}

		log.Debug("Request body decoded", slog.Any("request", req))

		// songDetail, err := GetSongInfoFromAPI(h, req)
		// if err != nil {
		// 	log.Error("Failed to get song info from external API", sl.Err(err))
		// 	c.JSON(http.StatusInternalServerError, resp.Error(err))
		// 	return
		// }
		// log.Debug("Retrieved song details from external API", slog.Any("request", songDetail))
		// songDetail.SongName = req.Song

		group := &models.Group{Name: req.Group}
		// createdGroup, err := h.songsUC.CreateGroup(n)
		// if err != nil {
		// 	log.Error("Failed to create group", sl.Err(err))
		// 	c.JSON(http.StatusBadRequest, err.Error())
		// 	return
		// }

		ReleaseDate, _ := time.Parse("16.07.2006", "16.07.2006")
		Text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
		Link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
		songDetail := models.Song{
			SongName:    req.Song,
			ReleaseDate: ReleaseDate,
			Text:        Text,
			Link:        Link,
		}

		// song := &models.Song{
		// 	GroupId:     createdGroup,
		// 	SongName:    req.Song,
		// 	ReleaseDate: songDetail.ReleaseDate,
		// 	Text:        songDetail.Text,
		// 	Link:        songDetail.Link,
		// }

		// createdSongs, err := h.songsUC.AddSong(song)
		// if err != nil {
		// 	log.Error("Failed to create users", sl.Err(err))
		// 	c.JSON(http.StatusBadRequest, err.Error())
		// 	return
		// }
		createdSongs, err := h.songsUC.AddGroupWithSongsTx(group, &songDetail)
		if err != nil {
			log.Error("Failed to add group and songs", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add group and songs"})
			return
		}

		c.JSON(http.StatusCreated, createdSongs)
	}
}

// // @Router /songs/create [post]
// func (h songsHandlers) AddSong() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		op := "songsHandlers.AddSong"
// 		log := h.log.With(slog.String("op", op))

// 		// Читаем запрос
// 		var req Request
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			log.Error("Failed to bind JSON", sl.Err(err))
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 			return
// 		}

// 		// Логика транзакции
// 		tx, err := h.db.Begin()
// 		if err != nil {
// 			log.Error("Failed to begin transaction", sl.Err(err))
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
// 			return
// 		}
// 		defer func() {
// 			if err != nil {
// 				tx.Rollback()
// 			} else {
// 				tx.Commit()
// 			}
// 		}()

// 		// log.Debug("Request body decoded", slog.Any("request", req))

// 		// songDetail, err := GetSongInfoFromAPI(h, req)
// 		// if err != nil {
// 		// 	log.Error("Failed to get song info from external API", sl.Err(err))
// 		// 	c.JSON(http.StatusInternalServerError, resp.Error(err))
// 		// 	return
// 		// }
// 		// log.Debug("Retrieved song details from external API", slog.Any("request", songDetail))

// 		n := &models.Group{Name: req.Group}
// 		createdGroup, err := h.songsUC.CreateGroup(n)
// 		if err != nil {
// 			log.Error("Failed to create group", sl.Err(err))
// 			c.JSON(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		ReleaseDate, _ := time.Parse("16.07.2006", "16.07.2006")
// 		Text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
// 		Link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
// 		songDetail := models.Song{
// 			ReleaseDate: ReleaseDate,
// 			Text:        Text,
// 			Link:        Link,
// 		}

// 		// Создаем песню
// 		song := &models.Song{
// 			GroupId:     createdGroup,
// 			SongName:    req.Song,
// 			ReleaseDate: songDetail.ReleaseDate, // Предполагается, что дата передается в запросе
// 			Text:        songDetail.Text,
// 			Link:        songDetail.Link,
// 		}

// 		createdSongs, err := h.songsUC.AddSong(tx, song)
// 		if err != nil {
// 			log.Error("Failed to add song", sl.Err(err))
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create song"})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, gin.H{"data": createdSongs})
// 	}
// }

// // @Router /songs/hello [get]
// func (h songsHandlers) Hello() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		op := "songsHandlers.Create"
// 		h.logger.Info("логгер работает!")
// 		h.logger.Error("Failed to get users", "err")
// 		c.JSON(http.StatusCreated, "Hello my friend")
// 	}
// }

func GetSongInfoFromAPI(h songsHandlers, req Request) (models.Song, error) {
	op := "songsHandlers.GetSongInfoFromAPI"

	log := h.log.With(
		slog.String("op", op),
	)
	group, song := req.Group, req.Song

	url := fmt.Sprintf("http://%s/info?group=%s&song=%s", h.cfg.Address, url.QueryEscape(group), url.QueryEscape(song))

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
