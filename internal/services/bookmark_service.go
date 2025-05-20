package services

import (
	"bookmarker/internal/models"
	"bookmarker/internal/repositories"
	"time"
)

// BookmarkService defines the service layer interface.
type BookmarkService interface {
	CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error)
	CreateBookmarkWithTags(url, title, description, thumbnail string, tags []string, createdAt time.Time) (models.Bookmark, error)
	GetBookmarkByID(id int) (models.Bookmark, error)
	GetBookmarkWithTags(id int) (models.Bookmark, error)
	ListBookmarks(page int, pageSize int) ([]models.Bookmark, error)
	ListBookmarksByTag(tagID int, page int, pageSize int) ([]models.Bookmark, error)
	ListBookmarksWithTags(page int, pageSize int) ([]models.Bookmark, error)
	UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error)
	UpdateBookmarkWithTags(id int, fields map[string]interface{}, tags []string) (models.Bookmark, error)
	// SearchBookmarks performs a paginated text search on title, url, or description
	SearchBookmarks(query string, page int, pageSize int) ([]models.Bookmark, error)
}

// bookmarkService implementation of the BookmarkService interface.
type bookmarkService struct {
	repo    repositories.BookmarkRepository
	tagRepo repositories.TagRepository
}

// NewBookmarkService creates a new instance of the bookmarkService.
func NewBookmarkService(repo repositories.BookmarkRepository) BookmarkService {
	return &bookmarkService{
		repo: repo,
	}
}

// NewBookmarkServiceWithTags creates a new instance of the bookmarkService with tagRepo.
func NewBookmarkServiceWithTags(repo repositories.BookmarkRepository, tagRepo repositories.TagRepository) BookmarkService {
	return &bookmarkService{
		repo:    repo,
		tagRepo: tagRepo,
	}
}

// CreateBookmark passes the bookmark to the repository for creation.
func (s *bookmarkService) CreateBookmark(url, title, description, thumbnail string) (models.Bookmark, error) {
	return s.repo.CreateBookmark(url, title, description, thumbnail, time.Now())
}

// CreateBookmarkWithTags creates a bookmark and associates tags.
func (s *bookmarkService) CreateBookmarkWithTags(url, title, description, thumbnail string, tags []string, createdAt time.Time) (models.Bookmark, error) {
	// Deduplicate tags
	tagSet := make(map[string]struct{})
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}
	uniqueTags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		uniqueTags = append(uniqueTags, tag)
	}

	// Create the bookmark
	bookmark, err := s.repo.CreateBookmark(url, title, description, thumbnail, createdAt)
	if err != nil {
		return bookmark, err
	}

	// Use new repo method to get/create tags and associate
	tagStructs, err := s.tagRepo.GetAndCreateTagsIfMissing(uniqueTags)
	if err != nil {
		return bookmark, err
	}
	for _, tag := range tagStructs {
		err := s.tagRepo.AddTagToBookmark(int(bookmark.ID), int(tag.ID))
		if err != nil {
			return bookmark, err
		}
	}
	return bookmark, nil
}

// GetBookmarkByID fetches a bookmark by its ID.
func (s *bookmarkService) GetBookmarkByID(id int) (models.Bookmark, error) {
	return s.repo.GetBookmarkByID(id)
}

// GetBookmarkWithTags fetches a bookmark by its ID and includes tags.
func (s *bookmarkService) GetBookmarkWithTags(id int) (models.Bookmark, error) {
	bookmark, err := s.repo.GetBookmarkByID(id)
	if err != nil {
		return bookmark, err
	}
	if s.tagRepo == nil {
		return bookmark, nil
	}
	tags, err := s.tagRepo.GetTagsForBookmark(int(bookmark.ID))
	if err != nil {
		return bookmark, err
	}
	bookmark.Tags = tags
	return bookmark, nil
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

// ListBookmarksWithTags retrieves paginated bookmarks and includes tags.
func (s *bookmarkService) ListBookmarksWithTags(page int, pageSize int) ([]models.Bookmark, error) {
	bookmarks, err := s.ListBookmarks(page, pageSize)
	if err != nil {
		return nil, err
	}
	if s.tagRepo == nil {
		return bookmarks, nil
	}
	for i := range bookmarks {
		tags, err := s.tagRepo.GetTagsForBookmark(int(bookmarks[i].ID))
		if err != nil {
			return nil, err
		}
		bookmarks[i].Tags = tags
	}
	return bookmarks, nil
}

// PatchBookmark updates only the provided fields of a bookmark.
func (s *bookmarkService) UpdateBookmark(id int, fields map[string]interface{}) (models.Bookmark, error) {
	return s.repo.UpdateBookmark(id, fields)
}

func (s *bookmarkService) UpdateBookmarkWithTags(id int, fields map[string]interface{}, tags []string) (models.Bookmark, error) {
	bookmark, err := s.repo.UpdateBookmark(id, fields)
	if err != nil {
		return bookmark, err
	}
	if s.tagRepo == nil {
		return bookmark, nil
	}

	// Deduplicate tags
	tagSet := make(map[string]struct{})
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}
	uniqueTags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		uniqueTags = append(uniqueTags, tag)
	}

	// Get or create tags
	tagStructs, err := s.tagRepo.GetAndCreateTagsIfMissing(uniqueTags)
	if err != nil {
		return bookmark, err
	}

	// Remove all existing tag associations for this bookmark
	// and add the new ones
	// This requires a new method in TagRepository: RemoveAllTagsFromBookmark
	if remover, ok := s.tagRepo.(interface {
		RemoveAllTagsFromBookmark(bookmarkID int) error
	}); ok {
		if err := remover.RemoveAllTagsFromBookmark(id); err != nil {
			return bookmark, err
		}
	} else {
		return bookmark, nil // or return an error if strict
	}

	// Add new tag associations
	for _, tag := range tagStructs {
		if err := s.tagRepo.AddTagToBookmark(id, int(tag.ID)); err != nil {
			return bookmark, err
		}
	}

	// Fetch updated tags
	bookmark.Tags, err = s.tagRepo.GetTagsForBookmark(id)
	if err != nil {
		return bookmark, err
	}
	return bookmark, nil
}

// SearchBookmarks performs a paginated text search on title, url, or description
func (s *bookmarkService) SearchBookmarks(query string, page int, pageSize int) ([]models.Bookmark, error) {
	offset := (page - 1) * pageSize
	return s.repo.SearchBookmarks(query, offset, pageSize)
}
