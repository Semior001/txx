package txx

import (
	"os"
	"testing"

	"context"
	"database/sql"

	"path"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRow struct {
	X int
	Y sql.NullString
}

func TestTxManager_Tx(t *testing.T) {
	db, txm := prepareTxManager(t)
	err := txm.Tx(context.Background(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context) error {
		_, err := db.ExecContext(ctx, `CREATE TABLE test (x INTEGER PRIMARY KEY, y TEXT)`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, `INSERT INTO test(x, y) VALUES(1, 'blah'), (2, NULL)`)
		return err
	})
	assert.NoError(t, err)

	rowscan, err := db.QueryContext(context.Background(), `SELECT * FROM test ORDER BY x`)
	require.NoError(t, err)
	var rows []testRow

	for rowscan.Next() {
		row := testRow{}
		err = rowscan.Scan(&row.X, &row.Y)
		require.NoError(t, err)
		rows = append(rows, row)
	}

	assert.Equal(t, []testRow{
		{X: 1, Y: sql.NullString{String: "blah", Valid: true}},
		{X: 2, Y: sql.NullString{}},
	}, rows)
}

func prepareTxManager(t *testing.T) (*sql.DB, *TxManager) {
	loc, err := os.MkdirTemp("", "txx_driver_sqlite")
	require.NoError(t, err, "failed to make temp dir")

	t.Cleanup(func() { assert.NoError(t, os.RemoveAll(loc)) })

	return NewTxManager(&sqlite3.SQLiteDriver{}, "file:"+path.Join(loc, "test.db"))
}
