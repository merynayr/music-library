package repository

import (
	"database/sql"
	"fmt"
	"music-library/internal/models"
	"music-library/internal/songs"
	"music-library/pkg/utils"
	"strings"
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
		addSongQuery,
		groupID,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(
		&s.ID,
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
	switch {
	case err == sql.ErrNoRows:
		err = tx.QueryRow(createGroupQuery, group.Name).Scan(&groupID)
		if err != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	case err != nil:
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

func (r *songRepo) GetSongs(pq *utils.PaginationQuery) (*models.SongsList, error) {
	const op = "song.repository.postgres.GetSongs"

	rows, err := r.db.Query(getSongsQuery, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songsList = make(map[string][]*models.Song, pq.GetSize())
	for rows.Next() {
		song := &models.Song{}
		group := &models.Group{}
		if err := rows.Scan(&song.ID, &group.Name, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songsList[group.Name] = append(songsList[group.Name], song)
	}

	return &models.SongsList{
		Page:  pq.GetPage(),
		Size:  pq.GetSize(),
		Songs: songsList,
	}, nil
}

func (r *songRepo) GetSongText(id string, pq *utils.PaginationQuery) (string, error) {
	const op = "song.repository.postgres.GetSongText"

	var text string
	row := r.db.QueryRow(getSongTextQuery, id)
	err := row.Scan(&text)
	switch {
	case err == sql.ErrNoRows:
		return "", fmt.Errorf("%s: %w", op, err)
	case err != nil:
		return "", fmt.Errorf("%s: %w", op, err)
	}

	verses := strings.Split(text, "\n\n")

	start := pq.GetOffset()
	if start >= len(verses) {
		errStr := fmt.Sprintf("No verses on page %d", pq.Page)
		err := fmt.Errorf("%s", errStr)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	end := start + pq.Size
	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := ""
	for i := start; i < end; i++ {
		paginatedVerses += verses[i]
	}
	return paginatedVerses, nil
}

func (r *songRepo) UpgradeGroupWithSongsTx(id string, Data map[string]interface{}) error {
	err := r.execTx(func(tx *sql.Tx) error {
		if value, exists := Data["group"]; exists {
			err := r.UpdateGroup(tx, id, value)
			if err != nil {
				return err
			}
		}

		err := r.UpdateSong(tx, id, Data)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (r *songRepo) UpdateSong(tx *sql.Tx, id string, Data map[string]interface{}) error {
	const op = "song.repository.postgres.UpdateSong"
	result, err := r.db.Exec(
		updateSongQuery,
		Data["song"],
		Data["releaseDate"],
		Data["text"],
		Data["link"],
		id,
	)

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

func (r *songRepo) UpdateGroup(tx *sql.Tx, id string, GroupID interface{}) error {
	const op = "song.repository.postgres.UpdateGroup"
	result, err := r.db.Exec(updateGroupQuery, GroupID, id)
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
