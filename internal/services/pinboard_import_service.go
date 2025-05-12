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
	"/src/assets/site1.png",
	"/src/assets/site2.png",
	"/src/assets/site3.png",
	"/src/assets/site4.png",
	"/src/assets/site5.png",
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
		
		_, err = s.BookmarkService.CreateBookmarkWithTags(
			pb.Href,
			pb.Description,
			description,
			thumbnail,
			tags,
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
