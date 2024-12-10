package songs

import (
	"music-library/internal/models"
)

// Songs use case
type UseCase interface {
	AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error)
	DeleteSong(id uint) error
	GetSongs() ([]models.Song, error)
	GetSongText(id string) (string, error)
	UpdateSong(id string, Data map[string]interface{}) error
}
