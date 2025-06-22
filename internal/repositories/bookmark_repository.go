package repositories

import (
	"bookmarker/internal/models"
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BookmarkRepository defines the interface for handling bookmarks with pagination support.
type BookmarkRepository interface {
	CreateBookmark(url, title, description, thumbnail string, createdAt time.Time) (models.Bookmark, error)
	GetBookmarkByID(id int) (models.Bookmark, error)
	ListBookmarks(offset int, limit int) ([]models.Bookmark, error)
	ListBookmarksByTag(tagID int, offset int, limit int) ([]models.Bookmark, error)
	UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error)
	// SearchBookmarks performs a paginated text search on title, url, or description
	SearchBookmarks(query string, offset int, limit int) ([]models.Bookmark, error)
	DeleteBookmark(id int) error // Add this method to the interface
}

type bookmarkRepository struct {
	db *pgxpool.Pool
}

// CreateBookmark adds a new bookmark to the database.
func (r bookmarkRepository) CreateBookmark(url, title, description, thumbnail string, createdAt time.Time) (models.Bookmark, error) {
	var bookmarkID int64
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO bookmarks (url, title, description, thumbnail, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		url, title, description, thumbnail, createdAt, createdAt,
	).Scan(&bookmarkID)
	if err != nil {
		return models.Bookmark{}, err
	}
	return models.Bookmark{
		ID:          bookmarkID,
		URL:         url,
		Title:       title,
		Description: &description,
		Thumbnail:   &thumbnail,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}, nil
}

// GetBookmarkByID retrieves a bookmark by its ID.
func (r bookmarkRepository) GetBookmarkByID(id int) (models.Bookmark, error) {
	query := `
		SELECT id, title, description, thumbnail, url, created_at, updated_at
		FROM bookmarks
		WHERE id = $1
	`
	row := r.db.QueryRow(context.Background(), query, id)
	var bookmark models.Bookmark
	err := row.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &bookmark.CreatedAt, &bookmark.UpdatedAt)
	if err != nil {
		return models.Bookmark{}, err
	}
	return bookmark, nil
}

// ListBookmarks retrieves a paginated list of bookmarks.
func (r bookmarkRepository) ListBookmarks(offset int, limit int) ([]models.Bookmark, error) {
	query := `
		SELECT id, title, description, thumbnail, url, created_at, updated_at
		FROM bookmarks
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookmarks []models.Bookmark
	for rows.Next() {
		var bookmark models.Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &bookmark.CreatedAt, &bookmark.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return bookmarks, nil
}

// ListBookmarksByTag retrieves a paginated list of bookmarks filtered by a tag.
func (r bookmarkRepository) ListBookmarksByTag(tagID int, offset int, limit int) ([]models.Bookmark, error) {
	query := `
		SELECT b.id, b.title, b.description, b.thumbnail, b.url, b.created_at, b.updated_at
		FROM bookmarks b
		INNER JOIN bookmarks_tags bt ON b.id = bt.bookmark_id
		WHERE bt.tag_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(context.Background(), query, tagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookmarks []models.Bookmark
	for rows.Next() {
		var bookmark models.Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &bookmark.CreatedAt, &bookmark.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return bookmarks, nil
}

// PatchBookmark updates only the provided fields and sets updated_at to now.
func (r bookmarkRepository) UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error) {
	if len(fields) == 0 {
		return r.GetBookmarkByID(id)
	}
	query := "UPDATE bookmarks SET "
	args := []interface{}{}
	i := 1
	for k, v := range fields {
		if i > 1 {
			query += ", "
		}
		query += k + " = $" + strconv.Itoa(i)
		args = append(args, v)
		i++
	}
	updatedAt := time.Now().UTC()
	query += ", updated_at = $" + strconv.Itoa(i) + " WHERE id = $" + strconv.Itoa(i+1)
	args = append(args, updatedAt, id)
	_, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return models.Bookmark{}, err
	}
	return r.GetBookmarkByID(id)
}

// SearchBookmarks performs a paginated text search on title, url, or description
func (r bookmarkRepository) SearchBookmarks(query string, offset int, limit int) ([]models.Bookmark, error) {
	likeQuery := "%" + query + "%"
	sqlQuery := `
		SELECT id, title, description, thumbnail, url, created_at, updated_at
		FROM bookmarks
		WHERE title ILIKE $1 OR url ILIKE $2 OR description ILIKE $3
		LIMIT $4 OFFSET $5
	`
	rows, err := r.db.Query(context.Background(), sqlQuery, likeQuery, likeQuery, likeQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookmarks []models.Bookmark
	for rows.Next() {
		var bookmark models.Bookmark
		err := rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &bookmark.CreatedAt, &bookmark.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return bookmarks, nil
}

// NewBookmarkRepository creates a new instance of bookmarkRepository.
func NewBookmarkRepository(db *pgxpool.Pool) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

// Implement DeleteBookmark for bookmarkRepository
func (r bookmarkRepository) DeleteBookmark(id int) error {
	// First, delete all tag relationships for this bookmark
	_, err := r.db.Exec(context.Background(), "DELETE FROM bookmarks_tags WHERE bookmark_id = $1", id)
	if err != nil {
		return err
	}
	// Then, delete the bookmark itself
	_, err = r.db.Exec(context.Background(), "DELETE FROM bookmarks WHERE id = $1", id)
	return err
}
