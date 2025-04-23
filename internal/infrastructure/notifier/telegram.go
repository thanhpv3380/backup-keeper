package notifier

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backup-keeper/internal/domain"
)

type TelegramNotifier struct {
	botToken string
	chatID   string
}

func NewTelegramNotifier(botToken string, chatID string) domain.Notifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
	}
}

func (n *TelegramNotifier) Notify(message string) error {
	log.Printf("Sending Telegram notification: %s", message)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)

	form := url.Values{}
	form.Add("chat_id", n.chatID)
	form.Add("text", message)

	// Tạo request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("failed to create request: ", err)
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed to send Telegram message: ", err)
		return fmt.Errorf("failed to send Telegram message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Println("telegram API error: ", string(body))
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	log.Println("Telegram notification sent successfully")
	return nil
}
