package sqltest

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// sqlite3Driver is the implementation of Driver for the "postgres" driver.
type sqlite3Driver struct{}

func (sqlite3Driver) Name() string {
	return "sqlite3"
}

func (sqlite3Driver) ParseDSN(dsn string) (DataSource, error) {
	if strings.Contains("dsn", ":memory:") {
		return nil, errors.New("in-memory databases are not supported")
	}

	dsn = strings.TrimPrefix(dsn, "file:")

	var params url.Values

	if pos := strings.IndexRune(dsn, '?'); pos != -1 {
		var err error
		params, err = url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return nil, err
		}

		dsn = dsn[:pos]
	}

	if dsn == "" {
		return nil, errors.New("DSN must contain a filename")
	}

	base := filepath.Base(dsn)
	ext := filepath.Ext(base)

	return &sqlite3DataSource{
		filepath.Dir(dsn),
		strings.TrimSuffix(base, ext),
		ext,
		params,
	}, nil
}

func (d sqlite3Driver) DataSourceForSQLite(filename string) (DataSource, error) {
	return d.ParseDSN("file:" + filename + "?mode=rwc")
}

// sqlite3DataSource is an implementation of DataSource for "sqlite3" driver.
type sqlite3DataSource struct {
	dir       string
	database  string
	extension string
	params    url.Values
}

func (ds sqlite3DataSource) DriverName() string {
	return "sqlite3"
}

func (ds sqlite3DataSource) DSN() string {
	return ds.filename() + "?" + ds.params.Encode()
}

func (ds sqlite3DataSource) DatabaseName() string {
	return ds.database
}

func (ds sqlite3DataSource) WithDatabaseName(database string) DataSource {
	ds.database = database
	return ds
}

func (ds sqlite3DataSource) Close() error {
	err := os.Remove(ds.filename())

	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func (ds sqlite3DataSource) filename() string {
	return filepath.Join(ds.dir, ds.database+ds.extension)
}
