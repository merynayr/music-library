package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultSize = 10
)

// Pagination query params
type PaginationQuery struct {
	Size int `json:"size,omitempty"`
	Page int `json:"page,omitempty"`
}

// Set page size
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("%s", "Size must be a positive value")
	}
	q.Size = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// Get OrderBy
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

// Get pagination query struct from
func GetPaginationFromCtx(c *gin.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.Query("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.Query("size")); err != nil {
		return nil, err
	}
	return q, nil
}
