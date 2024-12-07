package usecase

import (
	"music-library/config"
	"music-library/internal/models"
	"music-library/internal/songs"

	"github.com/sirupsen/logrus"
)

// News UseCase
type songsUC struct {
	cfg       *config.Config
	songsRepo songs.Repository
	logger    *logrus.Logger
}

// News UseCase constructor
func NewSongsUseCase(cfg *config.Config, songsRepo songs.Repository, logger *logrus.Logger) songs.UseCase {
	return &songsUC{cfg: cfg, songsRepo: songsRepo, logger: logger}
}

// Create song
func (s *songsUC) Create(songs *models.Song) (*models.Song, error) {
	return nil, nil
}
