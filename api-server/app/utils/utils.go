package utils

import (
	"database/sql"
	"errors"
	"os"
	"ride_sharing_api/app/simulator"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

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
	_, err := simulator.S.FsStat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func CreateDbFileIfNotExists(path string) error {
	_, err := simulator.S.FsStat(path)

	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := simulator.S.FsCreate(path)
	f.Close()

	return err
}

func SqlNullStr(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
