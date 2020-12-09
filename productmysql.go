package sqltest

import (
	"context"
	"database/sql"
)

// MySQLProtocol is an interface for drivers that use the MySQL wire protocol.
type MySQLProtocol interface {
	DataSourceForMySQL(
		user, pass,
		host, port,
		database string,
	) (DataSource, error)
}

// MySQLCompatibleProduct is a Product that is compatible with MySQL.
type MySQLCompatibleProduct struct {
	ProductName string
	DefaultPort string
}

// Name returns the human-readable name of the product.
func (p MySQLCompatibleProduct) Name() string {
	return p.ProductName
}

// DefaultDataSource returns the default data source to use to connect to the
// product.
func (p MySQLCompatibleProduct) DefaultDataSource(d Driver) (DataSource, error) {
	proto, ok := d.(MySQLProtocol)
	if !ok {
		return nil, ErrIncompatibleDriver
	}

	return proto.DataSourceForMySQL(
		"root", "rootpass",
		"127.0.0.1", p.DefaultPort,
		"mysql",
	)
}

// CreateDatabase creates a new database with the given name.
func (MySQLCompatibleProduct) CreateDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, "CREATE DATABASE `"+name+"`")
	return err
}

// DropDatabase drops the database with the given name.
func (MySQLCompatibleProduct) DropDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS `"+name+"`")
	return err
}
