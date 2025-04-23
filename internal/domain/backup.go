package domain

type BackupUseCase interface {
	Execute(dataSource string) error
}
