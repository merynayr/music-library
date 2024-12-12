package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type FilterQuery struct {
	Group       string `json:"group,omitempty"`
	SongName    string `json:"song,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

func GetFilterFromCtx(c *gin.Context) (*FilterQuery, error) {
	q := &FilterQuery{}
	if err := q.SetGroup(c.Query("group")); err != nil {
		return nil, err
	}
	if err := q.SetSong(c.Query("song")); err != nil {
		return nil, err
	}
	if query, flag := c.GetQuery("releaseDate"); flag {
		if err := q.SetReleaseDate(query); err != nil {
			return nil, err
		}
	} else {
		q.ReleaseDate = ""
	}

	if err := q.SetText(c.Query("text")); err != nil {
		return nil, err
	}
	if err := q.SetLink(c.Query("link")); err != nil {
		return nil, err
	}

	return q, nil
}

func (q *FilterQuery) SetGroup(groupQuery string) error {
	q.Group = groupQuery
	return nil
}

func (q *FilterQuery) SetSong(songQuery string) error {
	q.SongName = songQuery
	return nil
}

func (q *FilterQuery) SetReleaseDate(releaseDateQuery string) error {
	layouts := []string{
		"2006-01-02",
		"2006-01",
		"2006",
		"2006-01-02T15:04:05Z",
	}

	for _, layout := range layouts {
		_, err := time.Parse(layout, releaseDateQuery)
		if err == nil {
			q.ReleaseDate = "%" + releaseDateQuery + "%"
			return nil
		}
	}

	return fmt.Errorf("invalid date format: %s", releaseDateQuery)
}

func (q *FilterQuery) SetText(textQuery string) error {
	q.Text = textQuery
	return nil
}

func (q *FilterQuery) SetLink(linkQuery string) error {
	q.Link = linkQuery
	return nil
}
