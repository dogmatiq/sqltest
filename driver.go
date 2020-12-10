package sqltest

import (
	_ "github.com/go-sql-driver/mysql" // ensure all drivers are imported
	_ "github.com/jackc/pgx/v4/stdlib" // ensure all drivers are imported
	_ "github.com/lib/pq"              // ensure all drivers are imported
	_ "github.com/mattn/go-sqlite3"    // ensure all drivers are imported
)

// Driver is an interface for using a specific SQL driver with multiple database
// products.
type Driver interface {
	// Name returns the name of the driver, as passed to sql.Open().
	Name() string

	// IsAvailable returns true if this driver is available for use.
	IsAvailable() bool

	// ParseDSN parses a DSN. It returns an error if this DSN string is not
	// compatible with this driver.
	ParseDSN(dsn string) (DataSource, error)
}

// DataSource is a DSN tied to a specific driver.
type DataSource interface {
	// DriverName returns the name of the driver as used with sql.Open().
	DriverName() string

	// DSN returns the DSN string as used with sql.Open().
	DSN() string

	// DatabaseName returns the name of the database within the data source.
	DatabaseName() string

	// WithDatabaseName returns a clone of this data source that connects to a
	// different database.
	WithDatabaseName(name string) DataSource

	// Close releases any resources associated with the data source.
	Close() error
}

var (
	// MySQLDriver is the "mysql" driver (https://github.com/go-sql-driver/mysql).
	MySQLDriver Driver = mysqlDriver{}

	// PGXDriver is the "pgx" driver (https://github.com/jackc/pgx).
	PGXDriver Driver = pgxDriver{}

	// PostgresDriver is the "postgres" driver (https://github.com/lib/pq).
	PostgresDriver Driver = postgresDriver{}

	// SQLite3Driver is the "sqlite3" driver (github.com/mattn/go-sqlite3).
	SQLite3Driver Driver = sqlite3Driver{}

	// Drivers is a slice containing all known products.
	Drivers = []Driver{
		MySQLDriver,
		PGXDriver,
		PostgresDriver,
		SQLite3Driver,
	}
)
