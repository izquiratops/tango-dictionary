package utils

import (
	"os"
	"path/filepath"
)

func GetAbsolutePath(path string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	indexPath := filepath.Join(dir, path)

	return indexPath, err
}
