package repositories

import (
	"bookmarker/internal/models"
	"database/sql"
	"time"
)

type TagRepository interface {
	CreateTag(name string) (models.Tag, error)
	AddTagToBookmark(bookmarkID int, tagID int) error
	GetTagsForBookmark(bookmarkID int) ([]models.BookmarkTag, error)
	GetAndCreateTagsIfMissing(tagNames []string) ([]models.Tag, error)
	RemoveAllTagsFromBookmark(bookmarkID int) error
}

type tagRepository struct {
	db *sql.DB
}

// CreateTag adds a new tag to the database.
func (r tagRepository) CreateTag(name string) (models.Tag, error) {
	ts := time.Now().Unix()
	query := `INSERT INTO tags (name, created_at, updated_at) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, name, ts, ts)
	if (err != nil) {
		return models.Tag{}, err
	}

	tagID, err := result.LastInsertId()
	if (err != nil) {
		return models.Tag{}, err
	}

	// Fetch the created tag with timestamps
	var tag models.Tag
	tagQuery := `SELECT id, name, created_at, updated_at FROM tags WHERE id = ?`
	err = r.db.QueryRow(tagQuery, tagID).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if (err != nil) {
		return models.Tag{}, err
	}

	return tag, nil
}

// AddTagToBookmark associates a tag with a bookmark.
func (r tagRepository) AddTagToBookmark(bookmarkID int, tagID int) error {
	query := `INSERT INTO bookmarks_tags (bookmark_id, tag_id, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`

	_, err := r.db.Exec(query, bookmarkID, tagID)
	return err
}

// GetTagsForBookmark retrieves all tags associated with a bookmark.
func (r tagRepository) GetTagsForBookmark(bookmarkID int) ([]models.BookmarkTag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN bookmarks_tags bt ON t.id = bt.tag_id
		WHERE bt.bookmark_id = ?
	`

	rows, err := r.db.Query(query, bookmarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.BookmarkTag
	for rows.Next() {
		var tag models.BookmarkTag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetAndCreateTagsIfMissing accepts a slice of tag names, creates any missing tags, and returns all tag structs for the input names
func (r tagRepository) GetAndCreateTagsIfMissing(tagNames []string) ([]models.Tag, error) {
	if len(tagNames) == 0 {
		return nil, nil
	}

	// 1. Find which tags already exist
	query, args := "SELECT id, name, created_at, updated_at FROM tags WHERE name IN (", []interface{}{}
	for i, name := range tagNames {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, name)
	}
	query += ")"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existingTags := make(map[string]models.Tag)
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		existingTags[tag.Name] = tag
	}

	// 2. Find missing tags
	missing := []string{}
	for _, name := range tagNames {
		if _, found := existingTags[name]; !found {
			missing = append(missing, name)
		}
	}

	// 3. Create missing tags
	for _, name := range missing {
		_, err := r.CreateTag(name)
		if err != nil {
			return nil, err
		}
	}

	// 4. Query again to get all tag structs for input names
	rows2, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var tags []models.Tag
	for rows2.Next() {
		var tag models.Tag
		err := rows2.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// RemoveAllTagsFromBookmark removes all tags associated with a bookmark.
func (r tagRepository) RemoveAllTagsFromBookmark(bookmarkID int) error {
	query := `DELETE FROM bookmarks_tags WHERE bookmark_id = ?`
	_, err := r.db.Exec(query, bookmarkID)
	return err
}

// NewTagRepository creates a new instance of tagRepository.
func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}