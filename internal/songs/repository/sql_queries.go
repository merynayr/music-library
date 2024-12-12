package repository

const (
	addSongQuery = `INSERT INTO songs (group_id, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, song, release_date, text, link`

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
	  AND ($1 = '' OR groups.name ILIKE $1)
    AND ($2 = '' OR songs.song ILIKE $2) 
    AND ($3 = '' OR songs.release_date::text ILIKE $3)
    AND ($4 = '' OR songs.text ILIKE $4) 
    AND ($5 = '' OR songs.link ILIKE $5)
	ORDER BY group_name, song OFFSET $6 LIMIT $7
	`

	getSongTextQuery = `SELECT text FROM songs WHERE id = $1`

	updateSongQuery = `UPDATE songs 
	SET group_id = COALESCE(NULLIF($1, '')::integer, group_id),
		song = COALESCE(NULLIF($2, ''), song),
		release_date = COALESCE(NULLIF($3, '')::date, release_date), 
		text = COALESCE(NULLIF($4, ''), text), 
		link = COALESCE(NULLIF($5, ''), link)
	WHERE id = $6`
)
