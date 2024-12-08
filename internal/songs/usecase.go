package songs

import (
	"music-library/internal/models"
)

// Songs use case
type UseCase interface {
	AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error)
}
