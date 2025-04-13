package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backup-keeper/config"
	"backup-keeper/internal/infrastructure/collector"
	"backup-keeper/internal/infrastructure/notifier"
	"backup-keeper/internal/infrastructure/storage"
	"backup-keeper/internal/usecase"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB collector
	collector, err := collector.NewMongoDBCollector(
		cfg.MongoDB.URI,
		cfg.MongoDB.Database,
		cfg.MongoDB.Collections,
	)
	if err != nil {
		log.Fatalf("Failed to create MongoDB collector: %v", err)
	}
	defer collector.Close()

	// Initialize Telegram notifier
	notifier := notifier.NewTelegramNotifier(
		cfg.Telegram.BotToken,
		cfg.Telegram.ChatID,
	)

	// Initialize Google Drive storage
	storage, err := storage.NewGoogleDriveStorage(
		cfg.GoogleDrive.CredentialsFile,
		cfg.GoogleDrive.TokenFile,
	)
	if err != nil {
		log.Fatalf("Failed to create Google Drive storage: %v", err)
	}

	backupService := usecase.NewBackupUseCase(collector, storage, notifier)

	// Create cron scheduler
	c := cron.New(cron.WithLocation(cfg.Backup.TimeZone))

	// Add cron job
	_, err = c.AddFunc(cfg.Backup.CronSchedule, func() {
		log.Println("Starting scheduled backup...")
		backupService.Execute()
	})
	if err != nil {
		log.Fatalf("Failed to schedule backup job: %v", err)
	}

	// Run initial backup immediately
	go func() {
		log.Println("Running initial backup...")
		backupService.Execute()
	}()

	// Start cron scheduler
	c.Start()
	log.Printf("Scheduled backup configured to run at: %s (%s)",
		cfg.Backup.CronSchedule, cfg.Backup.TimeZone)

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop cron scheduler
	c.Stop()
	log.Println("Backup service stopped")
}
