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

	// DatabaseNameFromDSN returns the database name in the given DSN.
	DatabaseNameFromDSN(dsn string) (string, error)

	// DSNForSchemaManipulation returns a DSN that can be used to connect to a
	// database in order to create and drop other databases.
	DSNForSchemaManipulation(templateDSN string) (string, error)

	// DSNForTesting returns a DSN that connects to a specific database for
	// running tests.
	DSNForTesting(templateDSN, database string) (string, error)
}

var (
	// MySQLDriver is the "mysql" driver (https://github.com/go-sql-driver/mysql).
	MySQLDriver Driver = mysqlDriver{}

	// PGXDriver is the "pgx" driver (https://github.com/jackc/pgx).
	PGXDriver Driver = pgxDriver{}

	// PostgresDriver is the "postgres" driver (https://github.com/lib/pq).
	PostgresDriver Driver = postgresDriver{}

	// // SQLite3Driver is the "sqlite3" driver (https://github.com/mattn/go-sqlite3).
	// SQLite3Driver Driver = sqlite3Driver{}
)
