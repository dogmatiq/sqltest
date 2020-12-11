package sqltest

import (
	"context"
	"database/sql"
)

// PostgresProtocol is an interface for drivers that use the PostgreSQL wire
// protocol.
type PostgresProtocol interface {
	DataSourceForPostgres(
		user, pass,
		host, port,
		database string,
	) (DataSource, error)
}

// PostgresCompatibleProduct is a Product that is compatible with PostgreSQL.
type PostgresCompatibleProduct struct {
	ProductName string
	DefaultPort string
}

// Name returns the human-readable name of the product.
func (p PostgresCompatibleProduct) Name() string {
	return p.ProductName
}

// IsCompatibleWith return true if the product is compatible with d.
func (p PostgresCompatibleProduct) IsCompatibleWith(d Driver) bool {
	_, ok := d.(PostgresProtocol)
	return ok
}

// DefaultDataSource returns the default data source to use to connect to the
// product.
func (p PostgresCompatibleProduct) DefaultDataSource(d Driver) (DataSource, error) {
	return d.(PostgresProtocol).DataSourceForPostgres(
		"postgres", "rootpass",
		"127.0.0.1", p.DefaultPort,
		"", // default database
	)
}

// CreateDatabase creates a new database with the given name.
func (PostgresCompatibleProduct) CreateDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, `CREATE DATABASE "`+name+`"`)
	return err
}

// DropDatabase drops the database with the given name.
func (PostgresCompatibleProduct) DropDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, `DROP DATABASE IF EXISTS "`+name+`" WITH (FORCE)`)
	return err
}
