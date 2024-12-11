package models

import "time"

type Song struct {
	ID          int       `json:"id"`
	SongName    string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Group struct {
	Name string `json:"name"`
}

type SongsList struct {
	Page  int                `json:"page"`
	Size  int                `json:"size"`
	Songs map[string][]*Song `json:"songs"`
}
