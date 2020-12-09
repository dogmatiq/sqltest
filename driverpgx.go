package sqltest

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

// pgxDriver is the implementation of Driver for the "pgx" driver.
type pgxDriver struct{}

func (pgxDriver) Name() string {
	return "pgx"
}

func (pgxDriver) ParseDSN(dsn string) (DataSource, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return pgxDataSource{
		config: cfg,
		dsn:    dsn,
	}, nil
}

func (d pgxDriver) DataSourceForPostgres(
	user, pass,
	host, port,
	database string,
) (DataSource, error) {
	dsn := "sslmode=disable"

	if user != "" {
		dsn += " user=" + user
	}

	if pass != "" {
		dsn += " password=" + pass
	}

	if host != "" {
		dsn += " host=" + host
	}

	if port != "" {
		dsn += " port=" + port
	}

	if database != "" {
		dsn += " database=" + database
	}

	fmt.Println(dsn)

	return d.ParseDSN(dsn)
}

// pgxDataSource is an implementation of DataSource for "pgx" driver.
type pgxDataSource struct {
	config *pgx.ConnConfig

	dsn     string
	release bool
}

func (ds pgxDataSource) DriverName() string {
	return "pgx"
}

func (ds pgxDataSource) DSN() string {
	return ds.dsn
}

func (ds pgxDataSource) DatabaseName() string {
	return ds.config.Database
}

func (ds pgxDataSource) WithDatabaseName(database string) DataSource {
	cfg := ds.config.Copy()
	cfg.Database = database

	return pgxDataSource{
		config:  cfg,
		dsn:     stdlib.RegisterConnConfig(ds.config),
		release: true,
	}
}

func (ds pgxDataSource) Close() error {
	if ds.release {
		stdlib.UnregisterConnConfig(ds.dsn)
	}

	return nil
}
