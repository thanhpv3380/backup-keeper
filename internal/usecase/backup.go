package usecase

import (
	"backup-keeper/internal/domain"
	"backup-keeper/internal/utils"
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

func (uc *backupUseCase) Execute(dataSource string) error {
	data, err := uc.collector.Collect()
	if err != nil {
		uc.notifier.Notify("ℹ️ Data Source: " + dataSource + "\n\n❌ Backup failed during collection: " + err.Error())
		return err
	}

	filename := "backup_" + utils.GenerateTimestamp() + ".zip"
	if err := uc.storage.Save(filename, data); err != nil {
		uc.notifier.Notify("ℹ️ Data Source: " + dataSource + "\n\n❌ Backup failed during upload: " + err.Error())
		return err
	}

	uc.notifier.Notify("ℹ️ Data Source: " + dataSource + "\n\n✅ Backup successful: " + filename)
	return nil
}
