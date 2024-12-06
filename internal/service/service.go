package service

import (
	"music-library/internal/config"
	"music-library/internal/storage/songs"
)

type Service struct {
	cfg     *config.Config
	storage songs.Storage
}

func New(cfg *config.Config, storage songs.Storage) (*Service, error) {
	return &Service{
		cfg:     cfg,
		storage: storage}, nil
}

type Music interface {
	GetAllSongs() []string
	AddSong(name string) error
}

func (s *Service) GetAllSongs() []string {
	// Реализация метода
	return []string{"Song 1", "Song 2"}
}

func (s *Service) AddSong(name string) error {
	// Реализация метода
	return nil
}
