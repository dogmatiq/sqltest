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

func (postgresDriver) DSNForSchemaManipulation(templateDSN string) (string, error) {
	u, err := url.Parse(templateDSN)
	if err != nil {
		return "", err
	}

	u.Path = "" // connect to the default database

	return u.String(), nil
}

func (postgresDriver) DSNForTesting(templateDSN, database string) (string, error) {
	u, err := url.Parse(templateDSN)
	if err != nil {
		return "", err
	}

	u.Path = database

	return u.String(), nil
}
