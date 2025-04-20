package models

// Bookmark represents the bookmarks table in the database.
type Bookmark struct {
	ID          int64   		`json:"id"`
	Title       string  		`json:"title"`
	Description *string 		`json:"description,omitempty"`
	Thumbnail   *string 		`json:"thumbnail,omitempty"`
	URL         string  		`json:"url"`
	Tags        []BookmarkTag 	`json:"tags"`
	CreatedAt   int64   		`json:"created_at"`
	UpdatedAt   int64   		`json:"updated_at"`
}
