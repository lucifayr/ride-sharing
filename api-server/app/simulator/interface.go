package simulator

import (
	"database/sql"
	"io"
	"io/fs"
)

// !EXPERIMENTAL!
// Wrapper around all IO operations. `simulator/realworld.go` provides an
// implementation of this interface to use when running the application. Other
// implementations are used for simulation testing.
type Simulator interface {
	// Stat returns a FileInfo describing the named file.
	FsStat(name string) (fs.FileInfo, error)

	// Create creates or truncates the named file. If the file already exists,
	// it is truncated. If the file does not exist, it is created.
	FsCreate(name string) (File, error)

	// Log output for standard logger. Should be set on using
	// `log.SetOutput(Simulator.LogOutput())` on application startup.
	LogOutput() io.Writer

	// Open opens a database specified by its database driver name and a
	// driver-specific data source name, usually consisting of at least a
	// database name and connection information.
	SqlOpen(driverName string, dataSourceName string) (DB, error)
}

type DB interface {
	// Exec executes a query without returning any rows. The args are for any
	// placeholder parameters in the query.
	Exec(query string, args ...any) (sql.Result, error)
}

type File interface {
	// Closes the File, rendering it unusable for I/O.
	Close() error

	// implement if needed:
	// Chdir() error
	// Chmod(mode fs.FileMode) error
	// Chown(uid int, gid int) error
	// Fd() uintptr
	// Name() string
	// Read(b []byte) (n int, err error)
	// ReadAt(b []byte, off int64) (n int, err error)
	// ReadDir(n int) ([]fs.DirEntry, error)
	// ReadFrom(r io.Reader) (n int64, err error)
	// Readdir(n int) ([]fs.FileInfo, error)
	// Readdirnames(n int) (names []string, err error)
	// Seek(offset int64, whence int) (ret int64, err error)
	// SetDeadline(t time.Time) error
	// SetReadDeadline(t time.Time) error
	// SetWriteDeadline(t time.Time) error
	// Stat() (fs.FileInfo, error)
	// Sync() error
	// SyscallConn() (syscall.RawConn, error)
	// Truncate(size int64) error
	// Write(b []byte) (n int, err error)
	// WriteAt(b []byte, off int64) (n int, err error)
	// WriteString(s string) (n int, err error)
	// WriteTo(w io.Writer) (n int64, err error)
}