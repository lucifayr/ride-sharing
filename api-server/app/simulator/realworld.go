package simulator

import (
	"database/sql"
	"io"
	"io/fs"
	"os"
	"ride_sharing_api/app/assert"
)

type SimulatorRealWorld struct{}

type SFileRealWorld struct {
	inner *os.File
}

type SDBRealWorld struct {
	inner *sql.DB
}

func (s *SimulatorRealWorld) FsStat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (s *SimulatorRealWorld) FsCreate(name string) (File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return &SFileRealWorld{inner: f}, nil
}

func (s *SimulatorRealWorld) LogOutput() io.Writer {
	return os.Stdout
}

func (s *SimulatorRealWorld) SqlOpen(driverName string, dataSourceName string) (DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SDBRealWorld{inner: db}, nil
}

func (db *SDBRealWorld) Exec(query string, args ...any) (sql.Result, error) {
	assert.True(db != nil && db.inner != nil, "Trying to execute query on nil DB.", "query:", query)
	return db.inner.Exec(query, args...)
}

func (f *SFileRealWorld) Close() error {
	assert.True(f != nil && f.inner != nil, "Trying to close nil file.")
	return f.inner.Close()
}
