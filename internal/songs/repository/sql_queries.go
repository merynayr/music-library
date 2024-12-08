package repository

const (
	addSong = `INSERT INTO songs (group_id, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

	createGroup = `INSERT INTO groups (name) VALUES ($1) RETURNING id`

	checkExistGroup = `SELECT id FROM groups WHERE name = $1`
)
