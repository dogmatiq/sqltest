package sqlstub

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync/atomic"
)

// Connector is a test implementation of the driver.Connector interface.
type Connector struct {
	driver.Connector

	ConnectFunc func(context.Context) (driver.Conn, error)
	DriverFunc  func() driver.Driver
}

// Connect returns a connection to the database.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if c.ConnectFunc != nil {
		return c.ConnectFunc(ctx)
	}

	if c.Connector != nil {
		return c.Connector.Connect(ctx)
	}

	return nil, errors.New("sqlstub.Connector: no behavior configured")
}

// Driver returns the underlying Driver of the Connector.
func (c *Connector) Driver() driver.Driver {
	if c.DriverFunc != nil {
		return c.DriverFunc()
	}

	if c.Connector != nil {
		return c.Connector.Driver()
	}

	return &Driver{}
}

// Driver is a test implementation of the driver.Driver interface.
type Driver struct {
	driver.Driver

	OpenFunc func(string) (driver.Conn, error)
}

// Open returns a new connection to the database.
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	if d.OpenFunc != nil {
		return d.OpenFunc(dsn)
	}

	if d.Driver != nil {
		return d.Driver.Open(dsn)
	}

	return nil, errors.New("sqlstub.Driver: no behavior configured")
}

var counter uint32 // atomic

// RegisterDriver creates a Driver stub and registers it with the sql package.
//
// It returns the driver name that can be used with sql.Open() and the Driver
// stub itself.
func RegisterDriver() (string, *Driver) {
	n := fmt.Sprintf(
		"dogmatiq-sqltest-%d",
		atomic.AddUint32(&counter, 1),
	)

	d := &Driver{}
	sql.Register(n, d)

	return n, d
}
