package utils

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"path/filepath"
	embeddings "ride_sharing_api"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/database/migrations"
	"runtime"

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

func GetEnvRequired(key string) string {
	val, exists := os.LookupEnv(key)
	assert.True(exists, "Required environment variable isn't set.", "key:", key)
	return val
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !(errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission))
}

func CreateDbFileIfNotExists(path string) error {
	if FileExists(path) {
		return nil
	}

	f, err := os.Create(path)
	f.Close()

	return err
}

func SqlNullStr(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: *str, Valid: true}
}

func SqlNullStrWrapped(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}

func InitDb(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+dbFile)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	m := migrations.FromEmbedFs(embeddings.DbMigrations, "db/migrations")
	m.Up(db)

	return db, nil
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func ProjectRoot() string {
	// `basepath` is always relative to utils.go path
	return path.Join(basepath, "../..")
}
