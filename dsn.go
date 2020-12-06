package sqltest

import (
	"database/sql"
	"fmt"
	"sync"

	"go.uber.org/multierr"
)

// DSN is a DSN constructed by a driver. It may have associated resources that
// need to be cleaned up when the DSN is no longer required.
type DSN struct {
	Driver     string
	DataSource string
	Closer     func(*DSN) error

	m         sync.Mutex
	databases []*sql.DB
}

// Close releases any resources associated with the DSN.
func (d *DSN) Close() error {
	if d == nil {
		return nil
	}

	d.m.Lock()
	databases := d.databases
	d.databases = nil
	d.m.Unlock()

	var err error

	for _, db := range databases {
		err = multierr.Append(err, db.Close())
	}

	if d.Closer != nil {
		err = multierr.Append(err, d.Closer(d))
	}

	return err
}

// Open opens a database as per the parameters of this DSN.
//
// The database is closed when d.Close() is called.
func (d *DSN) Open() (*sql.DB, error) {
	db, err := sql.Open(d.Driver, d.DataSource)
	if err != nil {
		return nil, fmt.Errorf(
			"could not open temporary database (%s) using the '%s' driver: %w",
			d.DataSource,
			d.Driver,
			err,
		)
	}

	d.m.Lock()
	defer d.m.Unlock()

	d.databases = append(d.databases, db)

	return db, nil
}
