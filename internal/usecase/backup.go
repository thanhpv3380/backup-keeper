package usecase

import (
	"backup-keeper/internal/domain"
	"time"
)

type backupUseCase struct {
	collector domain.Collector
	storage   domain.Storage
	notifier  domain.Notifier
}

func NewBackupUseCase(
	collector domain.Collector,
	storage domain.Storage,
	notifier domain.Notifier,
) domain.BackupUseCase {
	return &backupUseCase{
		collector: collector,
		storage:   storage,
		notifier:  notifier,
	}
}

func (uc *backupUseCase) Execute() error {
	data, err := uc.collector.Collect()
	if err != nil {
		uc.notifier.Notify("❌ Backup failed during collection: " + err.Error())
		return err
	}

	filename := "backup_" + generateTimestamp() + ".json"
	if err := uc.storage.Save(filename, data); err != nil {
		uc.notifier.Notify("❌ Backup failed during upload: " + err.Error())
		return err
	}

	uc.notifier.Notify("✅ Backup successful: " + filename)
	return nil
}

func generateTimestamp() string {
	return time.Now().Format("20060102_150405")
}
