package sqltest

import (
	"context"
	"database/sql"
)

type mysqlCompatable struct{}

func (mysqlCompatable) CreateDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, "CREATE DATABASE `"+name+"`")
	return err
}

func (mysqlCompatable) DropDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS `"+name+"`")
	return err
}

type mysqlProduct struct {
	mysqlCompatable
}

func (p mysqlProduct) Name() string {
	return "MySQL"
}

func (p mysqlProduct) TemplateDSN(d Driver) (string, bool) {
	switch d {
	case MySQLDriver:
		return "root:rootpass@tcp(127.0.0.1:23306)/dogmatiq", true
	default:
		return "", false
	}
}

type mariaProduct struct {
	mysqlCompatable
}

func (p mariaProduct) Name() string {
	return "MariaDB"
}

func (p mariaProduct) TemplateDSN(d Driver) (string, bool) {
	switch d {
	case MySQLDriver:
		return "root:rootpass@tcp(127.0.0.1:23307)/dogmatiq", true
	default:
		return "", false
	}
}
