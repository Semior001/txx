package txx

import (
	"database/sql"
	"database/sql/driver"

	"context"

	"errors"

	"fmt"

	"github.com/Semior001/txx/proxy"
)

// ErrTxStarted indicates about the attempt of starting another transaction inside
// the transaction.
// TODO: implement savepoints
var ErrTxStarted = errors.New("transaction is already started")

// Runner aggregates all methods of the sqlx.DB and sqlx.Tx
type Runner interface {
	driver.ExecerContext
	driver.QueryerContext
	driver.ConnPrepareContext
}

// TxManager is used to make requests transactional
type TxManager struct {
	db  *sql.DB
	log Logger
}

// NewTxManager makes new instance of TxManager.
// TODO: add options - logger
func NewTxManager(drv driver.Driver, dsn string) (*sql.DB, *TxManager) {
	txm := &TxManager{log: nopLogger{}}

	conn := proxy.NewConnector(drv, dsn, proxy.Handlers{
		QueryContextHandler:   txm.queryContext,
		PrepareContextHandler: txm.prepareContext,
		BeginTxHandler:        txm.beginTx,
		ExecContextHandler:    txm.execContext,
	})

	db := sql.OpenDB(conn)
	return db, txm
}

// Tx runs the fn in the transaction.
func (m *TxManager) Tx(ctx context.Context, opts *sql.TxOptions, fn func(context.Context) error) error {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("grab conn: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			m.log.CloseConn(err)
		}
	}()

	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	err = conn.Raw(func(driverConn interface{}) error {
		dc, ok := driverConn.(Runner)
		if !ok {
			panic("driverConn passed from Raw is not a Runner")
		}

		if err = fn(putRunner(ctx, dc)); err != nil {
			return fmt.Errorf("fn returned err: %w", err)
		}

		return nil
	})
	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			m.log.Rollback(rberr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			m.log.Rollback(rberr)
		}
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (m *TxManager) beginTx(next driver.ConnBeginTx) driver.ConnBeginTx {
	return proxy.BeginTxFunc(
		func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
			if tx := getTx(ctx); tx != nil {
				return nil, ErrTxStarted
			}
			return next.BeginTx(ctx, opts)
		})
}

func (m *TxManager) queryContext(next driver.QueryerContext) driver.QueryerContext {
	return proxy.QueryContextFunc(
		func(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
			if tx := getTx(ctx); tx != nil {
				return tx.QueryContext(ctx, q, args)
			}
			return next.QueryContext(ctx, q, args)
		})
}

func (m *TxManager) prepareContext(next driver.ConnPrepareContext) driver.ConnPrepareContext {
	return proxy.PrepareContextFunc(
		func(ctx context.Context, q string) (driver.Stmt, error) {
			if tx := getTx(ctx); tx != nil {
				return tx.PrepareContext(ctx, q)
			}
			return next.PrepareContext(ctx, q)
		})
}

func (m *TxManager) execContext(next driver.ExecerContext) driver.ExecerContext {
	return proxy.ExecContextFunc(
		func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
			if tx := getTx(ctx); tx != nil {
				return tx.ExecContext(ctx, query, args)
			}
			return next.ExecContext(ctx, query, args)
		})
}
