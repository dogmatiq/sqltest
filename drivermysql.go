package sqltest

import (
	"net"

	"github.com/go-sql-driver/mysql"
)

// mysqlDriver is the implementation of Driver for the "mysql" driver.
type mysqlDriver struct{}

func (mysqlDriver) Name() string {
	return "mysql"
}

func (mysqlDriver) IsAvailable() bool {
	return true
}

func (mysqlDriver) ParseDSN(dsn string) (DataSource, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	return mysqlDataSource{cfg}, nil
}

func (d mysqlDriver) DataSourceForMySQL(
	user, pass,
	host, port,
	database string,
) (DataSource, error) {
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = pass
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(host, port)
	cfg.DBName = database

	return mysqlDataSource{cfg}, nil
}

// mysqlDataSource is an implementation of DataSource for "mysql" driver.
type mysqlDataSource struct {
	config *mysql.Config
}

func (ds mysqlDataSource) DriverName() string {
	return "mysql"
}

func (ds mysqlDataSource) DSN() string {
	return ds.config.FormatDSN()
}

func (ds mysqlDataSource) DatabaseName() string {
	return ds.config.DBName
}

func (ds mysqlDataSource) WithDatabaseName(database string) DataSource {
	cfg := *ds.config
	cfg.DBName = database
	return mysqlDataSource{&cfg}
}

func (ds mysqlDataSource) Close() error {
	return nil
}
