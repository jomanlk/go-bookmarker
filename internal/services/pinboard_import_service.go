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
	"https://1drv.ms/i/c/f2b9998cb2784732/IQTMsvabZ4tmSr-bD1f2Kfw4AfXistrI6TilBEGdDcbGN7w??width=660",
	"https://1drv.ms/i/c/f2b9998cb2784732/IQQcpBLxCcY_T64O2aoDxfRXAT5xp_c81LQaAe0qGnzk0Q4?width=660",
	"https://1drv.ms/i/c/f2b9998cb2784732/IQRZZxzdGVxLRLoFxNHn1yEbAar3TrdRUYIwApp7VfAWBk8?width=660",
	"https://1drv.ms/i/c/f2b9998cb2784732/IQQmzNd7DL0pQ4TATzp2h0yMASFnZMcUFXdLc_7U-oJEZQk?width=660",
	"https://1drv.ms/i/c/f2b9998cb2784732/IQR1lj5F4m37RJ9q_yC8b4_cAT3KdgW1D4Wjlmj6ZYB5MeM?width=660",
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

func parsePinboardTime(ts string) (int64, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
