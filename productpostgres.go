package sqltest

import (
	"context"
	"database/sql"
)

type postgresProduct struct{}

func (p postgresProduct) Name() string {
	return "PostgreSQL"
}

func (p postgresProduct) TemplateDSN(d Driver) (string, bool) {
	switch d {
	case PGXDriver:
		return "database=dogmatiq user=postgres password=rootpass sslmode=disable host=127.0.0.1 port=25432", true
	default:
		return "", false
	}
}

func (postgresProduct) CreateDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, `CREATE DATABASE "`+name+`"`)
	return err
}

func (postgresProduct) DropDatabase(
	ctx context.Context,
	db *sql.DB,
	name string,
) error {
	_, err := db.ExecContext(ctx, `DROP DATABASE IF EXISTS "`+name+`"`)
	return err
}
