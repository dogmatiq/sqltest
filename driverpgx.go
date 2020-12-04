package sqltest

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type pgxDriver struct{}

func (pgxDriver) Name() string {
	return "pgx"
}

func (pgxDriver) DatabaseNameFromDSN(dsn string) (string, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return "", err
	}

	return cfg.Database, nil
}

func (pgxDriver) DSNForSchemaManipulation(templateDSN string) (string, error) {
	cfg, err := pgx.ParseConfig(templateDSN)
	if err != nil {
		return "", err
	}

	cfg.Database = "" // connect to the default database

	return stdlib.RegisterConnConfig(cfg), nil
}

func (pgxDriver) DSNForTesting(templateDSN, database string) (string, error) {
	cfg, err := pgx.ParseConfig(templateDSN)
	if err != nil {
		return "", err
	}

	cfg.Database = database

	return stdlib.RegisterConnConfig(cfg), nil
}
