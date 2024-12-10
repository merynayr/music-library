package usecase

import (
	"log/slog"
	"music-library/config"
	"music-library/internal/models"
	"music-library/internal/songs"
)

// News UseCase
type songsUC struct {
	cfg       *config.Config
	songsRepo songs.Repository
	logger    *slog.Logger
}

// News UseCase constructor
func NewSongsUseCase(cfg *config.Config, songsRepo songs.Repository, logger *slog.Logger) songs.UseCase {
	return &songsUC{cfg: cfg, songsRepo: songsRepo, logger: logger}
}

func (s *songsUC) AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error) {
	song, err := s.songsRepo.AddGroupWithSongsTx(group, song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *songsUC) DeleteSong(id uint) error {
	err := s.songsRepo.DeleteSong(id)
	return err
}

func (s *songsUC) GetSongs() ([]models.Song, error) {
	songs, err := s.songsRepo.GetSongs()
	if err != nil {
		return nil, err
	}

	return songs, nil
}
