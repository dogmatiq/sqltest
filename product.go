package sqltest

import (
	"context"
	"database/sql"
)

// Product is an enumeration of the supported database products.
type Product interface {
	Name() string

	// TemplateDSN returns a "template" DSN to use for connecting to various
	// databases on the same server.
	TemplateDSN(d Driver) (string, bool)

	// CreateDatabase creates a new database with the given name.
	CreateDatabase(ctx context.Context, db *sql.DB, name string) error

	// DropDatabase drops the database with the given name.
	DropDatabase(ctx context.Context, db *sql.DB, name string) error
}

var (
	// MySQL is the Product for MySQL (https://www.mysql.com).
	MySQL Product = mysqlProduct{}

	// MariaDB is the Product for MariaDB (https://mariadb.org).
	MariaDB Product = mariaProduct{}

	// PostgreSQL is the Product for PostgreSQL (https://www.postgresql.org).
	PostgreSQL Product = postgresProduct{}
)
