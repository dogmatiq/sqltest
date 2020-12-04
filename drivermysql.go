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

func (mysqlDriver) DSNForSchemaManipulation(templateDSN string) (string, error) {
	// Mangle the database name in the parsed DSN so that we can connect to the
	// server using the "mysql" database.
	//
	// This is necessary simply because the Go MySQL driver can not connect
	// without choosing a particular database, and the "mysql" database is
	// guaranteed to exist.
	cfg, err := mysql.ParseDSN(templateDSN)
	if err != nil {
		return "", err
	}
	cfg.DBName = "mysql"

	return cfg.FormatDSN(), nil
}

func (mysqlDriver) DSNForTesting(templateDSN, database string) (string, error) {
	cfg, err := mysql.ParseDSN(templateDSN)
	if err != nil {
		return "", err
	}

	cfg.DBName = database

	return cfg.FormatDSN(), nil
}
