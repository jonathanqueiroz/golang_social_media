// models/post.go
package models

import "time"

type Post struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published time.Time `json:"published"`
	AuthorID  uint64    `json:"author_id"`
}
