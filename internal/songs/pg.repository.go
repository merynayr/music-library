package songs

import (
	"database/sql"
	"music-library/internal/models"
)

// Songs Repository
type Repository interface {
	AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error)
	AddSong(tx *sql.Tx, song *models.Song) (*models.Song, error)
	CreateGroup(tx *sql.Tx, group *models.Group) (int, error)
}
