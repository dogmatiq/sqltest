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

func (d pgxDriver) DSNForSchemaManipulation(templateDSN string) (*DSN, error) {
	cfg, err := pgx.ParseConfig(templateDSN)
	if err != nil {
		return nil, err
	}

	cfg.Database = "" // connect to the default database

	return &DSN{
		Driver:     d.Name(),
		DataSource: stdlib.RegisterConnConfig(cfg),
		Closer: func(dsn *DSN) error {
			stdlib.UnregisterConnConfig(dsn.DataSource)
			return nil
		},
	}, nil
}

func (d pgxDriver) DSNForTesting(templateDSN, database string) (*DSN, error) {
	cfg, err := pgx.ParseConfig(templateDSN)
	if err != nil {
		return nil, err
	}

	cfg.Database = database

	return &DSN{
		Driver:     d.Name(),
		DataSource: stdlib.RegisterConnConfig(cfg),
		Closer: func(dsn *DSN) error {
			stdlib.UnregisterConnConfig(dsn.DataSource)
			return nil
		},
	}, nil
}
