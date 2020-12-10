package sqltest

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

// postgresDriver is the implementation of Driver for the "postgres" driver.
type postgresDriver struct{}

func (postgresDriver) Name() string {
	return "postgres"
}

func (postgresDriver) IsAvailable() bool {
	return true
}

func (postgresDriver) ParseDSN(dsn string) (DataSource, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "postgres" {
		return nil, errors.New("unexpected URL scheme, must be 'postgres'")
	}

	if u.Port() == "" {
		u.Host = net.JoinHostPort(u.Host, "5432")
	}

	return postgresDataSource{u}, nil
}

func (d postgresDriver) DataSourceForPostgres(
	user, pass,
	host, port,
	database string,
) (DataSource, error) {
	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
	}

	q := url.Values{
		"sslmode": []string{"disable"},
	}

	u.RawQuery = q.Encode()

	if user != "" || pass != "" {
		u.User = url.UserPassword(user, pass)
	}

	if port != "" {
		u.Host = net.JoinHostPort(u.Host, port)
	}

	if database != "" {
		u.Path = "/" + database
	}

	return postgresDataSource{u}, nil
}

// postgresDataSource is an implementation of DataSource for "postgres" driver.
type postgresDataSource struct {
	url *url.URL
}

func (ds postgresDataSource) DriverName() string {
	return "postgres"
}

func (ds postgresDataSource) DSN() string {
	return ds.url.String()
}

func (ds postgresDataSource) DatabaseName() string {
	return strings.TrimPrefix(ds.url.Path, "/")
}

func (ds postgresDataSource) WithDatabaseName(database string) DataSource {
	u := *ds.url
	u.Path = "/" + database
	return postgresDataSource{&u}
}

func (ds postgresDataSource) Close() error {
	return nil
}
