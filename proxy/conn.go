package proxy

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// conn is a wrapper for connections to the real database
// allowing to call middlewares.
type conn struct {
	conn driver.Conn
	Handlers
}

// newConn makes new instance of conn.
func newConn(driverConn driver.Conn, handlers Handlers) *conn {
	handlers.fillStubs()
	return &conn{conn: driverConn, Handlers: handlers}
}

// QueryContext proxies QueryContext calls to the database, in order to call
// middlewares in between.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if qc, ok := c.conn.(driver.QueryerContext); ok {
		return c.QueryContextHandler(qc).QueryContext(ctx, query, args)
	}
	return nil, ErrDriverNotSupports("QueryContext")
}

// PrepareContext proxies PrepareContext calls to the database, in order to call
// middlewares in between.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if pc, ok := c.conn.(driver.ConnPrepareContext); ok {
		return c.PrepareContextHandler(pc).PrepareContext(ctx, query)
	}
	return nil, ErrDriverNotSupports("QueryContext")
}

// BeginTx proxies BeginTx calls to the database, in order to call middlewares
// in between.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if bt, ok := c.conn.(driver.ConnBeginTx); ok {
		return c.BeginTxHandler(bt).BeginTx(ctx, opts)
	}
	return nil, ErrDriverNotSupports("BeginTx")
}

// ExecContext proxies ExecContext calls to the database, in order to call
// middlewares in between.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if ec, ok := (c.conn).(driver.ExecerContext); ok {
		return c.ExecContextHandler(ec).ExecContext(ctx, query, args)
	}
	return nil, ErrDriverNotSupports("ExecContext")
}

// Close proxies Close calls to the database, in order to call middlewares in
// between.
func (c *conn) Close() error { return c.conn.Close() }

// Prepare proxies Prepare calls to the database, in order to call middlewares
// in between.
func (c *conn) Prepare(query string) (driver.Stmt, error) { return c.conn.Prepare(query) }

// Begin proxies Begin calls to the database, in order to call middlewares in
// between.
func (c *conn) Begin() (driver.Tx, error) { return c.conn.Begin() } //nolint:staticcheck // doesn't mean that we don't have to proxy it

// Ping proxies Ping calls to the database, in order to call middlewares in
// between.
func (c *conn) Ping(ctx context.Context) error {
	if pinger, ok := c.conn.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return ErrDriverNotSupports("Ping")
}

// CheckNamedValue proxies CheckNamedValue calls to the database, in order to
// call middlewares in between.
func (c *conn) CheckNamedValue(value *driver.NamedValue) error {
	if checker, ok := c.conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(value)
	}
	return ErrDriverNotSupports("CheckNamedValue")
}

// ResetSession proxies ResetSession calls to the database, in order to call
// middlewares in between.
func (c *conn) ResetSession(ctx context.Context) error {
	if resetter, ok := c.conn.(driver.SessionResetter); ok {
		return resetter.ResetSession(ctx)
	}
	return ErrDriverNotSupports("ResetSession")
}

// ErrDriverNotSupports indicates that the underlying driver doesn't support
// the called operation.
type ErrDriverNotSupports string

// Error returns string representation of the error.
func (e ErrDriverNotSupports) Error() string {
	return fmt.Sprintf("underlying driver doesn't support %q operation", string(e))
}
