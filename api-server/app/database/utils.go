package utils

import (
	"errors"
	"os"
)

func CreateDbFileIfNotExists(path string) error {
	_, err := os.Stat(path)

	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Create(path)
	f.Close()

	return err
}
