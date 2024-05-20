package models

import (
	"time"
)

type Notification struct {
	ID           uint64    `json:"id,omitempty"`
	UserID       uint64    `json:"user_id,omitempty"`
	Type         string    `json:"type,omitempty"`
	SourceUserID uint64    `json:"source_user_id,omitempty"`
	SourcePostID *uint64   `json:"source_post_id,omitempty"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	Name         string    `json:"name,omitempty"`
	Username     string    `json:"username,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	PostContent  *string   `json:"post_content,omitempty"`
	OthersTotal  int       `json:"others_total"`
}
