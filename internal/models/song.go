package models

import "time"

type Song struct {
	ID          int       `json:"id"`
	GroupId     int       `json:"group_id"`
	SongName    string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Group struct {
	Name string `json:"name"`
}
