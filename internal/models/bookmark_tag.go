package models

// BookmarkTag represents the bookmark_tag table in the database.
type BookmarkTag struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
}
