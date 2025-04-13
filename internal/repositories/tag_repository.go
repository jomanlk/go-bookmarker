package repositories

import (
	"bookmarker/internal/models"
	"database/sql"
)

type TagRepository interface {
	CreateTag(name string) (models.Tag, error)
	AddTagToBookmark(bookmarkID int, tagID int) error
	GetTagsForBookmark(bookmarkID int) ([]models.Tag, error)
}

type tagRepository struct {
	db *sql.DB
}

// CreateTag adds a new tag to the database.
func (r tagRepository) CreateTag(name string) (models.Tag, error) {
	query := `INSERT INTO tags (name) VALUES (?)`

	result, err := r.db.Exec(query, name)
	if err != nil {
		return models.Tag{}, err
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		return models.Tag{}, err
	}

	return models.Tag{
		ID:   tagID,
		Name: name,
	}, nil
}

// AddTagToBookmark associates a tag with a bookmark.
func (r tagRepository) AddTagToBookmark(bookmarkID int, tagID int) error {
	query := `INSERT INTO bookmark_tag (bookmark_id, tag_id) VALUES (?, ?)`

	_, err := r.db.Exec(query, bookmarkID, tagID)
	return err
}

// GetTagsForBookmark retrieves all tags associated with a bookmark.
func (r tagRepository) GetTagsForBookmark(bookmarkID int) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN bookmark_tag bt ON t.id = bt.tag_id
		WHERE bt.bookmark_id = ?
	`

	rows, err := r.db.Query(query, bookmarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
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

// NewTagRepository creates a new instance of tagRepository.
func NewTagRepository(db *sql.DB) TagRepository {
	return &tagRepository{db: db}
}