package sqltest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/multierr"
)

// Database is a database created for the purposes of testing.
type Database struct {
	Driver     Driver
	Product    Product
	DataSource DataSource

	m       sync.Mutex
	open    bool
	closers []func() error
}

// NewDatabase returns a new temporary test database for a specific pair of
// database product and SQL driver.
//
// The returned Database is an io.Closer that must be closed when the database
// is no longer needed. It is safe to close the Database even if this function
// returns an error.
func NewDatabase(
	ctx context.Context,
	d Driver,
	p Product,
) (_ *Database, err error) {
	baseDS, err := dataSource(d, p)
	if err != nil {
		return nil, err
	}
	defer baseDS.Close()

	testDS := baseDS.WithDatabaseName(
		generateTemporaryDatabaseName(),
	)

	multi, ok := p.(MultiDatabaseProduct)

	if !ok {
		return &Database{
			Driver:     d,
			Product:    p,
			DataSource: testDS,
			open:       true,
		}, nil
	}

	pool, err := openPool(p, baseDS)
	if err != nil {
		return nil, err
	}

	if err := multi.CreateDatabase(ctx, pool, testDS.DatabaseName()); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"unable to create a temporary %s database using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	return &Database{
		Driver:     d,
		Product:    p,
		DataSource: testDS,
		open:       true,
		closers: []func() error{
			pool.Close,
			func() error {
				// We assume that the context passed to NewDatabase() has
				// already been canceled at this point. Allow an additional 3
				// seconds to drop the temporary database.
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				return multi.DropDatabase(ctx, pool, testDS.DatabaseName())
			},
		},
	}, nil
}

// Open returns a database pool that connects to this database.
//
// The pool is closed when db.Close() is called.
func (db *Database) Open() (*sql.DB, error) {
	if db == nil {
		return nil, errors.New("attempted to open a database pool for a closed database")
	}

	db.m.Lock()
	defer db.m.Unlock()

	if !db.open {
		return nil, errors.New("attempted to open a database pool for a closed database")
	}

	pool, err := openPool(db.Product, db.DataSource)
	if err != nil {
		return nil, err
	}

	db.closers = append(db.closers, pool.Close)

	return pool, nil
}

// Close releases any resources associated with the DSN.
func (db *Database) Close() error {
	if db == nil {
		return nil
	}

	db.m.Lock()
	defer db.m.Unlock()

	if !db.open {
		return nil
	}

	db.open = false
	closers := db.closers
	db.closers = nil

	var err error

	for i := len(closers) - 1; i >= 0; i-- {
		err = multierr.Append(err, closers[i]())
	}

	if db.DataSource != nil {
		err = multierr.Append(err, db.DataSource.Close())
	}

	return err
}

// dataSource returns the data source to use for the given combination of driver
// and product.
//
// It first checks for an environment variable containing a DSN. If that is not
// present it askes the product to generate a default DSN.
func dataSource(d Driver, p Product) (DataSource, error) {
	if !p.IsCompatibleWith(d) {
		return nil, fmt.Errorf(
			"%s is incompatible with the '%s' driver",
			p.Name(),
			d.Name(),
		)
	}

	key := strings.ToUpper(fmt.Sprintf("DOGMATIQ_TEST_DSN_%s_%s", p.Name(), d.Name()))
	dsn := os.Getenv(key)

	if dsn != "" {
		ds, err := d.ParseDSN(dsn)
		if err != nil {
			return nil, fmt.Errorf(
				"can not parse the DSN in the %s environment variable: %w",
				key,
				err,
			)
		}

		return ds, nil
	}

	ds, err := p.DefaultDataSource(d)

	if err != nil {
		return nil, fmt.Errorf(
			"can not build a default %s DSN using the '%s' driver: %w",
			p.Name(),
			d.Name(),
			err,
		)
	}

	return ds, nil
}

var counter uint64 // atomic

// generateTemporaryDatabaseName returns a name for a temporary test database.
func generateTemporaryDatabaseName() string {
	return fmt.Sprintf(
		"test_%d_%d",
		os.Getpid(),
		atomic.AddUint64(&counter, 1),
	)
}

// openPool opens a database pool for the given data source.
func openPool(p Product, ds DataSource) (_ *sql.DB, err error) {
	pool, err := sql.Open(
		ds.DriverName(),
		ds.DSN(),
	)
	if err == nil {
		err = pool.Ping()

		if err != nil {
			pool.Close()
		}
	}

	if err != nil {
		return nil, fmt.Errorf(
			"unable to open a %s database pool using the '%s' driver: %w",
			p.Name(),
			ds.DriverName(),
			err,
		)
	}

	return pool, nil
}
