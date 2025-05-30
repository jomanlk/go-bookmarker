package controllers

import (
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TelegramController struct {
	DB *pgxpool.Pool
}

func NewTelegramController(db *pgxpool.Pool) *TelegramController {
	return &TelegramController{DB: db}
}

// TelegramWebhookHandler handles POST requests from the Telegram bot webhook
func (tc *TelegramController) TelegramWebhookHandler(c *gin.Context) {
	// Validate Telegram token
	if !validateTelegramToken(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or missing telegram token"})
		return
	}

	// Parse and log the incoming JSON payload
	update, ok := parseAndLogTelegramUpdate(c)
	if !ok {
		return
	}

	// Extract URL and tags from message
	url, tags, text, found := extractURLAndTagsFromMessage(update)
	if !found {
		c.JSON(http.StatusOK, gin.H{"status": "no valid url detected"})
		return
	}

	log.Printf("[TelegramWebhookHandler] URL detected: %q (tags: %v)", url, tags)
	// Initialize repositories and service
	bookmarkRepo := repositories.NewBookmarkRepository(tc.DB)
	tagRepo := repositories.NewTagRepository(tc.DB)
	bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)

	// Create the bookmark (title, description, thumbnail left empty)
	bookmark, err := bookmarkService.CreateBookmarkWithTags(url, "", "", "", tags, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark", "details": err.Error()})
		return
	}

	// Extract chat ID
	chatID := extractChatID(update)
	if chatID != 0 {
		sendTelegramConfirmation(chatID, bookmark)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "bookmark saved",
		"bookmark": bookmark,
	})
}

// validateTelegramToken checks the Telegram webhook token
func validateTelegramToken(c *gin.Context) bool {
	secretToken := os.Getenv("TELEGRAM_WEBHOOK_SECRET")
	token := c.Query("token")
	return token == secretToken && secretToken != ""
}

// parseAndLogTelegramUpdate parses and logs the incoming Telegram update
func parseAndLogTelegramUpdate(c *gin.Context) (map[string]interface{}, bool) {
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Telegram update"})
		return nil, false
	}
	payloadBytes, _ := json.MarshalIndent(update, "", "  ")
	log.Printf("[TelegramWebhookHandler] Received payload: %s", string(payloadBytes))
	return update, true
}

// extractURLAndTagsFromMessage extracts the URL and tags from the Telegram message
func extractURLAndTagsFromMessage(update map[string]interface{}) (string, []string, string, bool) {
	message, ok := update["message"].(map[string]interface{})
	if !ok {
		return "", nil, "", false
	}
	text, ok := message["text"].(string)
	if !ok || text == "" {
		return "", nil, "", false
	}
	var url string
	var tags []string
	words := make([]string, 0)
	for _, w := range splitBySpace(text) {
		if w != "" {
			words = append(words, w)
		}
	}
	if len(words) > 0 {
		url = words[0]
		if len(words) > 1 {
			tags = words[1:]
		}
	}
	if url == "" || !isValidURL(url) {
		log.Printf("[TelegramWebhookHandler] No valid URL detected in message: %q", text)
		return "", nil, text, false
	}
	return url, tags, text, true
}

// extractChatID extracts the chat ID from the Telegram update
func extractChatID(update map[string]interface{}) int64 {
	message, ok := update["message"].(map[string]interface{})
	if !ok {
		return 0
	}
	if chat, ok := message["chat"].(map[string]interface{}); ok {
		if id, ok := chat["id"].(float64); ok {
			return int64(id)
		}
	}
	return 0
}

// sendTelegramConfirmation sends a confirmation message to the user via Telegram
func sendTelegramConfirmation(chatID int64, bookmark *services.BookmarkWithTags) {
	telegramClient := services.NewTelegramApiClient()
	msg := "URL: " + bookmark.URL
	if len(bookmark.Tags) > 0 {
		tagNames := make([]string, len(bookmark.Tags))
		for i, tag := range bookmark.Tags {
			tagNames[i] = tag.Name
		}
		msg += "\nTags: " + joinStrings(tagNames, ", ")
	}
	if err := telegramClient.SendMessage(chatID, msg); err != nil {
		log.Printf("[TelegramWebhookHandler] Failed to send Telegram message: %v", err)
	}
}

// splitBySpace splits a string by spaces (helper for parsing)
func splitBySpace(s string) []string {
	result := []string{}
	start := 0
	for i, c := range s {
		if c == ' ' || c == '\t' || c == '\n' {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

// isValidURL validates if a string is a valid URL with http or https scheme
func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}

// joinStrings joins a slice of strings with a separator
func joinStrings(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for _, s := range elems[1:] {
		result += sep + s
	}
	return result
}
