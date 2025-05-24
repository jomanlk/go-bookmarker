package services

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

type PinboardBookmark struct {
	Href        string `json:"href"`
	Description string `json:"description"`
	Extended    string `json:"extended"`
	Time        string `json:"time"`
	Tags        string `json:"tags"`
}

var placeholderThumbnails = []string{
	"/placeholders/site1.png",
	"/placeholders/site2.png",
	"/placeholders/site3.png",
	"/placeholders/site4.png",
	"/placeholders/site5.png",
}

// PinboardImportService imports bookmarks from a Pinboard JSON export
// and uses the existing BookmarkService to create bookmarks with tags.
type PinboardImportService struct {
	BookmarkService BookmarkService
}

func NewPinboardImportService(bookmarkService BookmarkService) *PinboardImportService {
	return &PinboardImportService{BookmarkService: bookmarkService}
}

// ImportFromJSON reads Pinboard JSON from r and imports bookmarks.
func (s *PinboardImportService) ImportFromJSON(r io.Reader) error {
	var pinboardBookmarks []PinboardBookmark
	dec := json.NewDecoder(r)
	if err := dec.Decode(&pinboardBookmarks); err != nil {
		return err
	}
	if len(pinboardBookmarks) == 0 {
		return errors.New("no bookmarks found in pinboard export")
	}
	var err error
	for _, pb := range pinboardBookmarks {
		tags := parseTags(pb.Tags)
		thumbnail := randomThumbnail()
		description := pb.Extended

		// Parse the time field
		var createdAt time.Time
		if pb.Time != "" {
			createdAt, err = time.Parse(time.RFC3339, pb.Time)
			if err != nil {
				return err
			}
		} else {
			createdAt = time.Now()
		}

		_, err = s.BookmarkService.CreateBookmarkWithTags(
			pb.Href,
			pb.Description,
			description,
			thumbnail,
			tags,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Fields(tags)
	for i, t := range parts {
		parts[i] = strings.TrimSpace(t)
	}
	return parts
}

func randomThumbnail() string {
	rand.Seed(time.Now().UnixNano())
	return placeholderThumbnails[rand.Intn(len(placeholderThumbnails))]
}
