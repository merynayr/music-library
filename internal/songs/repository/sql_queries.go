package repository

const (
	addSongQuery = `INSERT INTO songs (group_id, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

	createGroupQuery = `INSERT INTO groups (name) VALUES ($1) RETURNING id`

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
	WHERE 1=1
	`

	getSongTextQuery = `SELECT text FROM songs WHERE id = $1`

	updateGroupQuery = `
        	UPDATE groups
        	SET name = $1
        	WHERE id = (SELECT group_id FROM songs WHERE id = $2)
    	`
	updateSongQuery = `UPDATE songs 
	SET song = COALESCE(NULLIF($1, ''), song),
		release_date = COALESCE(NULLIF($2, '')::date, release_date), 
		text = COALESCE(NULLIF($3, ''), text), 
		link = COALESCE(NULLIF($4, ''), link)
	WHERE id = $5`
)
