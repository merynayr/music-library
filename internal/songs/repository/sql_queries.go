package repository

const (
	addSong = `INSERT INTO songs (group_id, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

	createGroup = `INSERT INTO groups (name) VALUES ($1) RETURNING id`

	checkExistGroup = `SELECT id FROM groups WHERE name = $1`

	deleteSongQuery = `DELETE FROM songs WHERE id = $1`

	getSongsQuery = `SELECT 
		songs.id, 
		groups.name AS group_name, 
		songs.song AS song, 
		songs.release_date, 
		songs.text, 
		songs.link 
	FROM songs 
	JOIN groups ON songs.group_id = groups.id 
	WHERE 1=1`
)
