// models/post.go
package models

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID        uint64    `json:"id,omitempty"`
	ParentID  *uint64   `json:"parent_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	AuthorID  uint64    `json:"author_id,omitempty"`
}

func (post *Post) Prepare() error {
	if err := post.validate(); err != nil {
		return err
	}

	post.format()
	return nil
}

func (post *Post) validate() error {
	if post.Content == "" {
		return errors.New("content is required")
	}

	return nil
}

func (post *Post) format() {
	post.Content = strings.TrimSpace(post.Content)
}
