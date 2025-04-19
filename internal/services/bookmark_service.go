package services

import (
	"bookmarker/internal/models"
	"bookmarker/internal/repositories"
)

// BookmarkService defines the service layer interface.
type BookmarkService interface {
	CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error)
	GetBookmarkByID(id int) (models.Bookmark, error)
	ListBookmarks(page int, pageSize int) ([]models.Bookmark, error)
	ListBookmarksByTag(tagID int, page int, pageSize int) ([]models.Bookmark, error)
	UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error)
}

// bookmarkService implementation of the BookmarkService interface.
type bookmarkService struct {
	repo repositories.BookmarkRepository
}

// NewBookmarkService creates a new instance of the bookmarkService.
func NewBookmarkService(repo repositories.BookmarkRepository) BookmarkService {
	return &bookmarkService{
		repo: repo,
	}
}

// CreateBookmark passes the bookmark to the repository for creation.
func (s *bookmarkService) CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error) {
	return s.repo.CreateBookmark(url, title, description, thumbnail)
}

// GetBookmarkByID fetches a bookmark by its ID.
func (s *bookmarkService) GetBookmarkByID(id int) (models.Bookmark, error) {
	return s.repo.GetBookmarkByID(id)
}

// ListBookmarks retrieves paginated bookmarks.
func (s *bookmarkService) ListBookmarks(page int, pageSize int) ([]models.Bookmark, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListBookmarks(offset, pageSize)
}

// ListBookmarksByTag retrieves paginated bookmarks associated with a tag ID.
func (s *bookmarkService) ListBookmarksByTag(tagID int, page int, pageSize int) ([]models.Bookmark, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListBookmarksByTag(tagID, offset, pageSize)
}

 

// PatchBookmark updates only the provided fields of a bookmark.
func (s *bookmarkService) UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error) {
	return s.repo.UpdateBookmark(id, fields)
}
