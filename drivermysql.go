package sqltest

import (
	"github.com/go-sql-driver/mysql"
)

type mysqlDriver struct{}

func (mysqlDriver) Name() string {
	return "mysql"
}

func (mysqlDriver) DatabaseNameFromDSN(dsn string) (string, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", err
	}

	return cfg.DBName, nil
}

func (d mysqlDriver) DSNForSchemaManipulation(templateDSN string) (*DSN, error) {
	// Mangle the database name in the parsed DSN so that we can connect to the
	// server using the "mysql" database.
	//
	// This is necessary simply because the Go MySQL driver can not connect
	// without choosing a particular database, and the "mysql" database is
	// guaranteed to exist.
	cfg, err := mysql.ParseDSN(templateDSN)
	if err != nil {
		return nil, err
	}
	cfg.DBName = "mysql"

	return &DSN{
		Driver:     d.Name(),
		DataSource: cfg.FormatDSN(),
	}, nil
}

func (d mysqlDriver) DSNForTesting(templateDSN, database string) (*DSN, error) {
	cfg, err := mysql.ParseDSN(templateDSN)
	if err != nil {
		return nil, err
	}

	cfg.DBName = database

	return &DSN{
		Driver:     d.Name(),
		DataSource: cfg.FormatDSN(),
	}, nil
}
