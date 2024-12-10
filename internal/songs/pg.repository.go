package songs

import (
	"database/sql"
	"music-library/internal/models"
)

// Songs Repository
type Repository interface {
	AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error)
	AddSong(tx *sql.Tx, groupID int, song *models.Song) (*models.Song, error)
	CreateGroup(tx *sql.Tx, group *models.Group) (int, error)

	DeleteSong(id uint) error

	GetSongs() ([]models.Song, error)
	GetSongText(id string) (string, error)

	UpgradeGroupWithSongsTx(id string, Data map[string]interface{}) error
	UpdateSong(tx *sql.Tx, id string, Data map[string]interface{}) error
	UpdateGroup(tx *sql.Tx, id string, GroupName interface{}) error
}
