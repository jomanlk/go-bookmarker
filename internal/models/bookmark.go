package models

import "time"

// Bookmark represents the bookmarks table in the database.
type Bookmark struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description *string     `json:"description,omitempty"`
	Thumbnail   *string     `json:"thumbnail,omitempty"`
	URL         string      `json:"url"`
	Tags        []BookmarkTag `json:"tags"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
