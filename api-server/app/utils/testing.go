package utils

import (
	"database/sql"
	"os"
	"path"
	"ride_sharing_api/app/assert"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

var existingDBs = make([]string, 0)

func InitTestDB(setupFilePath string) *sql.DB {
	assert.False(slices.Contains(existingDBs, setupFilePath), "Test database was initialized more than once!", "path:", setupFilePath)
	existingDBs = append(existingDBs, setupFilePath)

	dbName := strings.TrimSuffix(path.Base(setupFilePath), path.Ext(setupFilePath))
	dbDir := path.Join(ProjectRoot(), "db/testing/instances")
	dbPath := path.Join(dbDir, dbName+".db")

	prefix := strings.Split(dbName, "-")[0]
	_, err := strconv.Atoi(prefix)
	assert.Nil(err, "Invalid prefix for setup file. Expected valid number (e.g. 0123-). Received", prefix)

	syscall.Umask(0)
	err = os.MkdirAll(dbDir, 0777)
	f, err := os.Create(dbPath)
	assert.Nil(err)
	f.Close()

	db, err := InitDb(dbPath)
	assert.Nil(err)

	ddlBytes, err := os.ReadFile(setupFilePath)
	assert.Nil(err)
	ddl := string(ddlBytes)

	for _, line := range strings.Split(ddl, "\n") {
		requirePath, found := strings.CutPrefix(line, "-- :require ")
		if !found {
			break
		}

		requiredDdl, err := os.ReadFile(path.Join(path.Dir(setupFilePath), requirePath))
		assert.Nil(err)

		_, err = db.Exec(string(requiredDdl))
		assert.Nil(err)
	}

	_, err = db.Exec(ddl)
	assert.Nil(err)

	return db
}
