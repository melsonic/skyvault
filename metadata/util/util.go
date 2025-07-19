package util

import (
	"errors"
	"fmt"
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

func FormatHashedChunks(hashedChunks []string) string {

	hashes := "{"
	for i := range hashedChunks {
		hashes += fmt.Sprintf(`"%s"`, hashedChunks[i])
		if i == len(hashedChunks)-1 {
			break
		}
		hashes += ","
	}
	hashes += "}"
	
	return hashes
}