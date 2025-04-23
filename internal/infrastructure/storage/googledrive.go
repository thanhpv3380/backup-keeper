package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"backup-keeper/internal/domain"
	"backup-keeper/internal/utils"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GoogleDriveStorage struct {
	service  *drive.Service
	folderId string
}

func NewGoogleDriveStorage(credentialsFile string, folderId string) (domain.Storage, error) {
	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &GoogleDriveStorage{service: srv, folderId: folderId}, nil
}

func (s *GoogleDriveStorage) Save(fileName string, data interface{}) error {
	filePath, ok := data.(string)
	if !ok {
		return fmt.Errorf("invalid data type, expected string")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer func() {
		file.Close()
		if err := utils.DeleteFile(filePath); err != nil {
			log.Println("⚠️ Warning: Failed to delete file upload: " + filePath + " - error: " + err.Error())
		}
	}()

	fileMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{s.folderId},
	}

	req := s.service.Files.Create(fileMetadata).
		Media(file, googleapi.ChunkSize(10*1024*1024)). // 10MB per chunk
		Context(context.Background())

	_, err = req.Do()
	if err != nil {
		return fmt.Errorf("failed to upload file using resumable upload: %v", err)
	}

	log.Printf("Uploaded %s to Google Drive folder %s", fileName, s.folderId)
	return nil
}
