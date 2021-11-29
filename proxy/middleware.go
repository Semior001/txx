package proxy

import (
	"context"
	"database/sql/driver"
)

// QueryContextMiddleware is used to make a chain of calls before reaching the
// real database.
type QueryContextMiddleware func(next driver.QueryerContext) driver.QueryerContext

// QueryContextFunc is an adapter to use ordinary functions as driver.QueryerContext.
type QueryContextFunc func(context.Context, string, []driver.NamedValue) (driver.Rows, error)

// QueryContext calls the wrapped function.
func (f QueryContextFunc) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return f(ctx, query, args)
}

// ExecContextMiddleware is used to make a chain of calls before reaching the
// real database.
type ExecContextMiddleware func(next driver.ExecerContext) driver.ExecerContext

// ExecContextFunc is an adapter to use ordinary functions as driver.ExecerContext.
type ExecContextFunc func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)

// ExecContext calls the wrapped function.
func (f ExecContextFunc) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return f(ctx, query, args)
}

// PrepareContextMiddleware is used to make a chain of calls before reaching the
// real database.
type PrepareContextMiddleware func(next driver.ConnPrepareContext) driver.ConnPrepareContext

// PrepareContextFunc is an adapter to use ordinary functions as driver.ConnPrepareContext.
type PrepareContextFunc func(ctx context.Context, query string) (driver.Stmt, error)

// PrepareContext calls the wrapped function.
func (f PrepareContextFunc) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return f(ctx, query)
}

// BeginTxMiddleware is used to make a chain of calls before reaching the
// real database.
type BeginTxMiddleware func(next driver.ConnBeginTx) driver.ConnBeginTx

// BeginTxFunc is an adapter to use ordinary functions as driver.ConnBeginTx.
type BeginTxFunc func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error)

// BeginTx calls the wrapped function.
func (f BeginTxFunc) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return f(ctx, opts)
}
