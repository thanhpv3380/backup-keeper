package config

import (
	"time"
)

type Config struct {
	MongoDB struct {
		URI         string
		Database    string
		Collections []string
	}

	Telegram struct {
		BotToken string
		ChatID   int64
	}

	GoogleDrive struct {
		CredentialsFile string
		TokenFile       string
	}

	Backup struct {
		CronSchedule string
		TimeZone     *time.Location
	}
}

func Load() *Config {
	cfg := &Config{}

	// MongoDB config
	cfg.MongoDB.URI = "mongodb://localhost:27017"
	cfg.MongoDB.Database = "myapp"
	cfg.MongoDB.Collections = []string{"users", "products", "orders"}

	// Telegram config
	cfg.Telegram.BotToken = "YOUR_TELEGRAM_BOT_TOKEN"
	cfg.Telegram.ChatID = 123456789 // Your chat ID

	// Google Drive config
	cfg.GoogleDrive.CredentialsFile = "credentials.json"
	cfg.GoogleDrive.TokenFile = "YOUR_GOOGLE_DRIVE_FOLDER_ID"

	// Backup schedule (every day at midnight)
	cfg.Backup.CronSchedule = "0 0 * * *"
	cfg.Backup.TimeZone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")

	return cfg
}
