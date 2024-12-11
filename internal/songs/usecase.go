package songs

import (
	"music-library/internal/models"
	"music-library/pkg/utils"
)

// Songs use case
type UseCase interface {
	AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error)
	DeleteSong(id uint) error
	GetSongs(pq *utils.PaginationQuery) (*models.SongsList, error)
	GetSongText(id string, pq *utils.PaginationQuery) (string, error)
	UpdateSong(id string, Data map[string]interface{}) error
}
