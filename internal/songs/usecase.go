package songs

import (
	"music-library/internal/models"
)

// Songs use case
type UseCase interface {
	Create(songs *models.Song) (*models.Song, error)
}
