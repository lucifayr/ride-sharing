package simulator

import (
	"context"
	"crypto/rand"
	"database/sql"
	"io"
	"io/fs"
	"net/http"
	"os"
	"ride_sharing_api/app/assert"
)

type SimulatorRealWorld struct{}

type FileRealWorld struct {
	inner *os.File
}

type DBRealWorld struct {
	inner *sql.DB
}

type HTTPMuxRealWorld struct {
	inner *http.ServeMux
}

func (s *SimulatorRealWorld) FsStat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (s *SimulatorRealWorld) FsCreate(name string) (File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return &FileRealWorld{inner: f}, nil
}

func (s *SimulatorRealWorld) NewHttpServerMux() HTTPMux {
	return &HTTPMuxRealWorld{inner: http.NewServeMux()}
}

func (s *SimulatorRealWorld) LogOutput() io.Writer {
	return os.Stdout
}

func (s *SimulatorRealWorld) SqlOpen(driverName string, dataSourceName string) (DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DBRealWorld{inner: db}, nil
}

func (s *SimulatorRealWorld) RandCrypto(b []byte) {
	rand.Read(b)
}

func (s *SimulatorRealWorld) HttpGet(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func (db *DBRealWorld) Exec(query string, args ...any) (sql.Result, error) {
	assert.True(db != nil && db.inner != nil, "Trying to call 'Exec' on nil DB.", "query:", query)
	return db.inner.Exec(query, args...)
}

func (db *DBRealWorld) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	assert.True(db != nil && db.inner != nil, "Trying to call 'ExecContext' on nil DB.", "query:", query)
	return db.inner.ExecContext(ctx, query, args...)
}

func (db *DBRealWorld) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	assert.True(db != nil && db.inner != nil, "Trying to call 'PrepareContext' on nil DB.", "query:", query)
	return db.inner.PrepareContext(ctx, query)
}

func (db *DBRealWorld) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	assert.True(db != nil && db.inner != nil, "Trying to call 'QueryContext' on nil DB.", "query:", query)
	return db.inner.QueryContext(ctx, query, args...)
}

func (db *DBRealWorld) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	assert.True(db != nil && db.inner != nil, "Trying to call 'QueryRowContext' on nil DB.", "query:", query)
	return db.inner.QueryRowContext(ctx, query, args...)
}

func (m *HTTPMuxRealWorld) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	assert.True(m != nil && m.inner != nil, "Trying to set handler function on nil mux handler.")
	m.inner.HandleFunc(pattern, handler)
}

func (m *HTTPMuxRealWorld) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assert.True(m != nil && m.inner != nil, "Trying to serve HTTP on nil mux handler.")
	m.inner.ServeHTTP(w, r)
}

func (f *FileRealWorld) Close() error {
	assert.True(f != nil && f.inner != nil, "Trying to close nil file.")
	return f.inner.Close()
}
