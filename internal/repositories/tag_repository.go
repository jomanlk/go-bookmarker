package repositories

import (
	"bookmarker/internal/models"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository interface {
	CreateTag(name string) (models.Tag, error)
	AddTagToBookmark(bookmarkID int, tagID int) error
	GetTagsForBookmark(bookmarkID int) ([]models.BookmarkTag, error)
	GetAndCreateTagsIfMissing(tagNames []string) ([]models.Tag, error)
	GetTagByName(name string) (models.Tag, error)
	RemoveAllTagsFromBookmark(bookmarkID int) error
	ListAllTags() ([]models.Tag, error)
	ListTags(page int, limit int) ([]models.Tag, error)
}

type tagRepository struct {
	db *pgxpool.Pool
}

// CreateTag adds a new tag to the database.
func (r tagRepository) CreateTag(name string) (models.Tag, error) {
	ts := time.Now().UTC()
	var tagID int64
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO tags (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`,
		name, ts, ts,
	).Scan(&tagID)
	if err != nil {
		return models.Tag{}, err
	}
	var tag models.Tag
	err = r.db.QueryRow(context.Background(),
		`SELECT id, name, created_at, updated_at FROM tags WHERE id = $1`, tagID,
	).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return models.Tag{}, err
	}
	return tag, nil
}

// AddTagToBookmark associates a tag with a bookmark.
func (r tagRepository) AddTagToBookmark(bookmarkID int, tagID int) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO bookmarks_tags (bookmark_id, tag_id, created_at) VALUES ($1, $2, NOW())`,
		bookmarkID, tagID,
	)
	return err
}

// GetTagsForBookmark retrieves all tags associated with a bookmark.
func (r tagRepository) GetTagsForBookmark(bookmarkID int) ([]models.BookmarkTag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN bookmarks_tags bt ON t.id = bt.tag_id
		WHERE bt.bookmark_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, bookmarkID)
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
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tags, nil
}

func (r tagRepository) GetTagByName(name string) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name, created_at, updated_at FROM tags WHERE name = $1`, name,
	).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return models.Tag{}, err
	}
	return tag, nil
}

// GetAndCreateTagsIfMissing accepts a slice of tag names, creates any missing tags, and returns all tag structs for the input names
func (r tagRepository) GetAndCreateTagsIfMissing(tagNames []string) ([]models.Tag, error) {
	if len(tagNames) == 0 {
		return nil, nil
	}
	// 1. Find which tags already exist
	placeholders := make([]string, len(tagNames))
	args := make([]interface{}, len(tagNames))
	for i, name := range tagNames {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args[i] = name
	}
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name IN (` + strings.Join(placeholders, ",") + `)`
	rows, err := r.db.Query(context.Background(), query, args...)
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
	missing := []string{}
	for _, name := range tagNames {
		if _, found := existingTags[name]; !found {
			missing = append(missing, name)
		}
	}
	for _, name := range missing {
		_, err := r.CreateTag(name)
		if err != nil {
			return nil, err
		}
	}
	// Query again to get all tag structs for input names
	rows2, err := r.db.Query(context.Background(), query, args...)
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
	_, err := r.db.Exec(context.Background(), `DELETE FROM bookmarks_tags WHERE bookmark_id = $1`, bookmarkID)
	return err
}

// ListAllTags retrieves all tags from the database (no pagination)
func (r tagRepository) ListAllTags() ([]models.Tag, error) {
	rows, err := r.db.Query(context.Background(), `SELECT id, name, created_at, updated_at FROM tags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tags, nil
}

// ListTags retrieves paginated tags from the database
func (r tagRepository) ListTags(page int, limit int) ([]models.Tag, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(context.Background(), `SELECT id, name, created_at, updated_at FROM tags LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tags, nil
}

// NewTagRepository creates a new instance of tagRepository.
func NewTagRepository(db *pgxpool.Pool) TagRepository {
	return &tagRepository{db: db}
}