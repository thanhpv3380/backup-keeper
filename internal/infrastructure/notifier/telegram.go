package notifier

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"backup-keeper/internal/domain"
)

type TelegramNotifier struct {
	botToken string
	chatID   int64
}

func NewTelegramNotifier(botToken string, chatID int64) domain.Notifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
	}
}

func (n *TelegramNotifier) Notify(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)

	form := url.Values{}
	form.Add("chat_id", fmt.Sprintf("%d", n.chatID))
	form.Add("text", message)

	resp, err := http.PostForm(apiURL, form)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}
