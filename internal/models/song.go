package models

import "time"

type Song struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group"`
	SongName    string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Group struct {
	Name string `json:"name"`
}
