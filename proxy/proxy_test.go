//go:build cgo

package proxy

import (
	"testing"

	"os"

	"path"

	"database/sql"

	"context"
	"database/sql/driver"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	queryCtx := context.WithValue(context.Background(), "t", "test")

	queryCalled := false
	execCalls := 0

	db := prepareDB(t, Handlers{
		QueryContextHandler: func(next driver.QueryerContext) driver.QueryerContext {
			return QueryContextFunc(func(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
				queryCalled = true
				assert.Equal(t, queryCtx, ctx)
				return next.QueryContext(ctx, q, args)
			})
		},
		ExecContextHandler: func(next driver.ExecerContext) driver.ExecerContext {
			return ExecContextFunc(func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
				execCalls++
				return next.ExecContext(ctx, query, args)
			})
		},
	})

	_, err := db.ExecContext(context.Background(), `CREATE TABLE test (x INTEGER PRIMARY KEY, y TEXT)`)
	assert.NoError(t, err)
	assert.Equal(t, 1, execCalls)

	type testrow struct {
		X int
		Y sql.NullString
	}

	_, err = db.ExecContext(context.Background(), `INSERT INTO test(x, y) VALUES(1, 'blah'), (2, NULL)`)
	assert.NoError(t, err)
	assert.Equal(t, 2, execCalls)

	rowscan, err := db.QueryContext(queryCtx, `SELECT * FROM test ORDER BY x`)
	assert.NoError(t, err)
	assert.True(t, queryCalled)

	var rows []testrow

	for rowscan.Next() {
		row := testrow{}
		err = rowscan.Scan(&row.X, &row.Y)
		require.NoError(t, err)
		rows = append(rows, row)
	}

	assert.Equal(t, []testrow{
		{X: 1, Y: sql.NullString{String: "blah", Valid: true}},
		{X: 2, Y: sql.NullString{}},
	}, rows)
}

func prepareDB(t *testing.T, handlers Handlers) *sql.DB {
	loc, err := os.MkdirTemp("", "txx_driver_sqlite")
	require.NoError(t, err, "failed to make temp dir")

	t.Cleanup(func() { assert.NoError(t, os.RemoveAll(loc)) })

	cn := NewConnector(&sqlite3.SQLiteDriver{}, "file:"+path.Join(loc, "test.db"), handlers)
	return sql.OpenDB(cn)
}
