package models

import "time"

// BookmarkTag represents the bookmark_tag table in the database.
type BookmarkTag struct {
	BookmarkID int64     `json:"bookmark_id"`
	TagID      int64     `json:"tag_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
