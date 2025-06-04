package util

import (
	"errors"
	"path/filepath"
)

type fileExtension string

func GetFileExtension(fileName string) (fileExtension, error) {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return "", errors.New("unsupported file name")
	}
	return fileExtension(ext), nil
}
