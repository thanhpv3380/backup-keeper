package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func WriteBatchToJson(batch []bson.M, filePath string) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, doc := range batch {
		if err := encoder.Encode(doc); err != nil {
			return "", err
		}
	}
	return filePath, nil
}

func ZipFiles(filePaths []string, outputZip string) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, path := range filePaths {
		if err := addFileToZip(zipWriter, path); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
