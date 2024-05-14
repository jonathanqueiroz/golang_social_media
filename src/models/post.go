// models/post.go
package models

import (
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID               uint64    `json:"id,omitempty"`
	ParentID         *uint64   `json:"parent_id,omitempty"`
	Content          string    `json:"content,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	AuthorID         uint64    `json:"author_id,omitempty"`
	AuthorName       string    `json:"author_name,omitempty"`
	Username         string    `json:"username,omitempty"`
	TotalLikes       uint64    `json:"total_likes"`
	CurrentUserLiked bool      `json:"current_user_liked"`
	TotalReplies     uint64    `json:"total_replies"`
	Replies          []Post    `json:"replies,omitempty"`
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

	if post.AuthorID == 0 {
		return errors.New("author is required")
	}

	return nil
}

func (post *Post) format() {
	post.Content = strings.TrimSpace(post.Content)
}
