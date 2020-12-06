package sqltest

import (
	"net/url"
	"strings"
)

type postgresDriver struct{}

func (postgresDriver) Name() string {
	return "postgres"
}

func (postgresDriver) DatabaseNameFromDSN(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(u.Path, "/"), nil
}

func (d postgresDriver) DSNForSchemaManipulation(templateDSN string) (*DSN, error) {
	u, err := url.Parse(templateDSN)
	if err != nil {
		return nil, err
	}

	u.Path = "" // connect to the default database

	return &DSN{
		Driver:     d.Name(),
		DataSource: u.String(),
	}, nil
}

func (d postgresDriver) DSNForTesting(templateDSN, database string) (*DSN, error) {
	u, err := url.Parse(templateDSN)
	if err != nil {
		return nil, err
	}

	u.Path = database

	return &DSN{
		Driver:     d.Name(),
		DataSource: u.String(),
	}, nil
}
