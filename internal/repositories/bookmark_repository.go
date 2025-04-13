package repositories

import (
	"bookmarker/internal/models"
	"database/sql"
	"time"
)

// BookmarkRepository defines the interface for handling bookmarks with pagination support.
type BookmarkRepository interface {
	CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error)
	GetBookmarkByID(id int) (models.Bookmark, error)
	ListBookmarks(offset int, limit int) ([]models.Bookmark, error)
	ListBookmarksByTag(tagID int, offset int, limit int) ([]models.Bookmark, error)
}

type bookmarkRepository struct {
	db *sql.DB
}
 

// CreateBookmark adds a new bookmark to the database.
func (r bookmarkRepository) CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error) {
    tx, err := r.db.Begin()
    if err != nil {
        return models.Bookmark{}, err
    }
    defer tx.Rollback()

    createdAt := time.Now().Unix()
    result, err := tx.Exec(`
                INSERT INTO bookmarks (url, title, description, thumbnail, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?)
        `, url, title, description, thumbnail, createdAt, createdAt)
    if err != nil {
        return models.Bookmark{}, err
    }

    bookmarkID, err := result.LastInsertId()
    if err != nil {
        return models.Bookmark{}, err
    }

    err = tx.Commit()
    if err != nil {
        return models.Bookmark{}, err
    }

    return models.Bookmark{
        ID:          bookmarkID,
        URL:         url,
        Title:       title,
        Description: &description,
        Thumbnail:   &thumbnail,
        CreatedAt:   time.Unix(createdAt, 0),
        UpdatedAt:   time.Unix(createdAt, 0),
    }, nil
}

// GetBookmarkByID retrieves a bookmark by its ID.
func (r bookmarkRepository) GetBookmarkByID(id int) (models.Bookmark, error) {
	query := `
		SELECT id, title, description, thumbnail, url, created_at, updated_at
		FROM bookmarks
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	var bookmark models.Bookmark
	var createdAt, updatedAt int64

	err := row.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return models.Bookmark{}, nil // No result found
	}
	if err != nil {
		return models.Bookmark{}, err
	}

	bookmark.CreatedAt = time.Unix(createdAt, 0)
	bookmark.UpdatedAt = time.Unix(updatedAt, 0)
	return bookmark, nil
}

// ListBookmarks retrieves a paginated list of bookmarks.
func (r bookmarkRepository) ListBookmarks(offset int, limit int) ([]models.Bookmark, error) {
    query := `
        SELECT id, title, description, thumbnail, url, created_at, updated_at
        FROM bookmarks
        LIMIT ? OFFSET ?
    `

    rows, err := r.db.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var bookmarks []models.Bookmark
    for rows.Next() {
        var bookmark models.Bookmark
        var createdAt, updatedAt int64 // Temporary variables for timestamps

        err := rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Description, &bookmark.Thumbnail, &bookmark.URL, &createdAt, &updatedAt)
        if err != nil {
            return nil, err
        }

        // Convert UNIX timestamps to time.Time
        bookmark.CreatedAt = time.Unix(createdAt, 0)
        bookmark.UpdatedAt = time.Unix(updatedAt, 0)

        bookmarks = append(bookmarks, bookmark)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    // Return an empty slice if no bookmarks are found
    if len(bookmarks) == 0 {
        return []models.Bookmark{}, nil
    }

    return bookmarks, nil
}

// ListBookmarksByTag retrieves a paginated list of bookmarks filtered by a tag.
func (r bookmarkRepository) ListBookmarksByTag(tagID int, offset int, limit int) ([]models.Bookmark, error) {
	query := `
		SELECT b.id, b.title, b.description, b.thumbnail, b.url, b.created_at, b.updated_at
		FROM bookmarks b
		INNER JOIN bookmark_tag bt ON b.id = bt.bookmark_id
		WHERE bt.tag_id = ?
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, tagID, limit, offset)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return an empty slice if no bookmarks are found
    if len(bookmarks) == 0 {
        return []models.Bookmark{}, nil
    }

	return bookmarks, nil
}

// NewBookmarkRepository creates a new instance of bookmarkRepository.
func NewBookmarkRepository(db *sql.DB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}
