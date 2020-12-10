package sqltest

import (
	"context"
	"database/sql"
)

// Product is a specific database product such as MySQL or MariaDB.
//
// The product correlates with a running service that tests are run against.
// Many products share a common wire protocol or are even entirely compatible.
type Product interface {
	// Name returns the human-readable name of the product.
	Name() string

	// IsCompatibleWith return true if the product is compatible with d.
	IsCompatibleWith(d Driver) bool

	// DefaultDataSource returns the default data source to use to connect to
	// the product.
	//
	// The returned data source must contain enough information to connect to
	// the product when running under the sqltest Docker stack or under a CI
	// workflow.
	//
	// d is the driver that is being used to connect to the product. If the
	// driver can not be used to connect to this product it returns
	// ErrIncompatibleDriver.
	DefaultDataSource(d Driver) (DataSource, error)
}

// MultiDatabaseProduct is a product that supports multiple databases on the
// same "server" or endpoint.
type MultiDatabaseProduct interface {
	Product

	// CreateDatabase creates a new database with the given name.
	CreateDatabase(ctx context.Context, db *sql.DB, name string) error

	// DropDatabase drops the database with the given name.
	DropDatabase(ctx context.Context, db *sql.DB, name string) error
}

var (
	// MySQL is the Product for MySQL (https://www.mysql.com).
	MySQL Product = MySQLCompatibleProduct{
		ProductName: "MySQL",
		DefaultPort: "23306",
	}

	// MariaDB is the Product for MariaDB (https://mariadb.org).
	MariaDB Product = MySQLCompatibleProduct{
		ProductName: "MariaDB",
		DefaultPort: "23307",
	}

	// PostgreSQL is the Product for PostgreSQL (https://www.postgresql.org).
	PostgreSQL Product = PostgresCompatibleProduct{
		ProductName: "PostgreSQL",
		DefaultPort: "25432",
	}

	// SQLite is the Product for SQLite (https://www.sqlite.org).
	SQLite Product = sqliteProduct{}

	// Products is a slice containing all known products.
	Products = []Product{
		MySQL,
		MariaDB,
		PostgreSQL,
		SQLite,
	}
)
