package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type MongoDB struct {
	URI      string
	Database string
}

type Telegram struct {
	BotToken string
	ChatID   string
}

type GoogleDrive struct {
	CredentialsFile string
	FolderId        string
}

type Backup struct {
	DataSource   string
	CronSchedule string
	TimeZone     string
}

type Config struct {
	MongoDB     MongoDB
	Telegram    Telegram
	GoogleDrive GoogleDrive
	Backup      Backup
}

var (
	Cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found. Using default values or system environment variables.")
		}

		Cfg = &Config{
			MongoDB: MongoDB{
				URI:      getEnv("MONGODB_URI", ""),
				Database: getEnv("MONGODB_DATABASE", "jcce"),
			},
			Telegram: Telegram{
				BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
				ChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
			},
			GoogleDrive: GoogleDrive{
				CredentialsFile: getEnv("GOOGLE_DRIVE_CREDENTIALS_FILE", ""),
				FolderId:        getEnv("GOOGLE_DRIVE_FOLDER_ID", ""),
			},
			Backup: Backup{
				DataSource:   getEnv("BACKUP_DATA_SOURCE", ""),
				CronSchedule: getEnv("BACKUP_CRON_SCHEDULE", ""),
				TimeZone:     getEnv("BACKUP_TIMEZONE", "Asia/Ho_Chi_Minh"),
			},
		}
	})

	return Cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
