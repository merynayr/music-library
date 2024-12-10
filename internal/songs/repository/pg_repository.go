package repository

import (
	"database/sql"
	"fmt"
	"music-library/internal/models"
	"music-library/internal/songs"
)

// News Repository
type songRepo struct {
	db *sql.DB
}

// News repository constructor
func NewSongsRepository(db *sql.DB) songs.Repository {
	return &songRepo{
		db: db,
	}
}

func (r *songRepo) execTx(fn func(*sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (r *songRepo) AddGroupWithSongsTx(group *models.Group, song *models.Song) (*models.Song, error) {
	var result *models.Song
	err := r.execTx(func(tx *sql.Tx) error {
		var err error

		groupID, err := r.CreateGroup(tx, group)
		if err != nil {
			return err
		}

		result, err = r.AddSong(tx, groupID, song)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func (r *songRepo) AddSong(tx *sql.Tx, groupID int, song *models.Song) (*models.Song, error) {
	const op = "song.repository.postgres.AddSong"

	var s models.Song
	err := tx.QueryRow(
		addSong,
		groupID,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(
		&s.ID,
		&s.GroupName,
		&s.SongName,
		&s.ReleaseDate,
		&s.Text,
		&s.Link)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &s, nil
}

func (r *songRepo) CreateGroup(tx *sql.Tx, group *models.Group) (int, error) {
	const op = "song.repository.postgres.CreateGroup"

	var groupID int
	err := r.db.QueryRow(checkExistGroup, group.Name).Scan(&groupID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow(createGroup, group.Name).Scan(&groupID)
		if err != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	} else if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return groupID, nil
}

func (r *songRepo) DeleteSong(id uint) error {
	const op = "song.repository.postgres.DeleteSong"

	result, err := r.db.Exec(deleteSongQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *songRepo) GetSongs() ([]models.Song, error) {
	const op = "song.repository.postgres.GetSongs"

	var args []interface{}
	rows, err := r.db.Query(getSongsQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}
