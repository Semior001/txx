package proxy

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// Connector is a wrapper for nested driver to pass
// it to sql.OpenDB in order to wrap all connections to the real
// database into a proxy.
type Connector struct {
	driver   driver.Driver
	dsn      string
	handlers Handlers
}

// NewConnector makes new instance of Connector.
func NewConnector(drv driver.Driver, dsn string, handlers Handlers) *Connector {
	return &Connector{driver: drv, dsn: dsn, handlers: handlers}
}

// Connect tries to call wrapped driver.DriverContext.OpenConnector, and, if fails,
// calls the Open of the wrapped driver.
func (dc *Connector) Connect(ctx context.Context) (conn driver.Conn, err error) {
	if oc, ok := dc.driver.(driver.DriverContext); ok {
		connector, err := oc.OpenConnector(dc.dsn)
		if err != nil {
			return nil, fmt.Errorf("open connector: %w", err)
		}
		if conn, err = connector.Connect(ctx); err != nil {
			return nil, err
		}

		return conn, nil
	}

	if conn, err = dc.driver.Open(dc.dsn); err != nil {
		return nil, err
	}

	return newConn(conn, dc.handlers), nil
}

// Driver returns the wrapped driver
func (dc *Connector) Driver() driver.Driver {
	return dc.driver
}
