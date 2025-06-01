package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TelegramApiClient struct {
	BotToken string
}

func NewTelegramApiClient() *TelegramApiClient {
	return &TelegramApiClient{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

// SendMessage sends a message to a Telegram chat
func (c *TelegramApiClient) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.BotToken)
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram sendMessage failed: %s", resp.Status)
	}
	return nil
}