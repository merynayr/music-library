package repository

import (
	"context"
	"database/sql"
	"music-library/internal/models"
	"music-library/internal/songs"
)

// News Repository
type songRepo struct {
	db *sql.DB
}

// News repository constructor
func NewSongsRepository(db *sql.DB) songs.Repository {
	return &songRepo{db: db}
}

func (r *songRepo) Create(ctx context.Context, songs *models.Song) (*models.Song, error) {
	return nil, nil
}
