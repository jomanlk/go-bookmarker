package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// URLPreviewResponse holds the preview information for a URL
// (renamed from LinkPreviewResponse)
type URLPreviewResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}

// URLPreviewApiClient is a client for the LinkPreview API
// (renamed from LinkPreviewApiClient)
type URLPreviewApiClient struct {
	ApiKey string
}

// NewURLPreviewApiClient creates a new URLPreviewApiClient
// (renamed from NewLinkPreviewApiClient)
func NewURLPreviewApiClient() *URLPreviewApiClient {
	return &URLPreviewApiClient{
		ApiKey: os.Getenv("LINK_PREVIEW_API_KEY"),
	}
}

// Fetch fetches preview information for a given URL
// (renamed receiver and response type)
func (c *URLPreviewApiClient) Fetch(url string) (*URLPreviewResponse, error) {
	apiUrl := fmt.Sprintf("https://api.linkpreview.net/?q=%s", url)
	request, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Linkpreview-Api-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linkpreview API error: %s", resp.Status)
	}

	var result URLPreviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
