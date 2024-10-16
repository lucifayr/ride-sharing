package utils

import (
	"errors"
	"os"
)

// Get the element at index	`idx` of `slice`. Returns `nil, false` if the index
// is out of bounds.
func SliceGet[T any](slice []T, idx int) (*T, bool) {
	if idx >= len(slice) {
		return nil, false
	}

	return &slice[idx], true
}

func IdxOf[T comparable](slice []T, predicate func(item T) bool) int {
	for idx, value := range slice {
		if predicate(value) {
			return idx
		}
	}

	return -1
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func CreateDbFileIfNotExists(path string) error {
	_, err := os.Stat(path)

	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Create(path)
	f.Close()

	return err
}
