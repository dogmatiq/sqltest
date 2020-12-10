package sqltest

import (
	"os"
	"path/filepath"
)

// SQLiteProtocol is an interface for drivers that use SQLite 3 database files.
type SQLiteProtocol interface {
	DataSourceForSQLite(filename string) (DataSource, error)
}

// sqliteProduct is the SQLite product.
type sqliteProduct struct{}

func (sqliteProduct) Name() string {
	return "SQLite"
}

func (sqliteProduct) IsCompatibleWith(d Driver) bool {
	_, ok := d.(SQLiteProtocol)
	return ok
}

func (sqliteProduct) DefaultDataSource(d Driver) (DataSource, error) {
	return d.(SQLiteProtocol).DataSourceForSQLite(
		filepath.Join(os.TempDir(), "dogmatiq.sqlite3"),
	)
}
