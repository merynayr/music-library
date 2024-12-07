package songs

import (
	"context"
	"music-library/internal/models"
)

// Songs Repository
type Repository interface {
	Create(ctx context.Context, songs *models.Song) (*models.Song, error)
}
