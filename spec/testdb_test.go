package spec

import (
	"context"
	"database/sql"
	"io"
	"strings"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabenv"
)

type mockDBConnection struct{}

func (m *mockDBConnection) Get(f func(db krabdb.DB) error) error {
	db, err := krabdb.Connect(krabenv.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close()

	return f(&testDB{recorder: &strings.Builder{}, db: db})
}

type testDB struct {
	db       *sqlx.DB
	recorder io.StringWriter
}

func (d *testDB) GetDatabase() *sqlx.DB {
	return d.db
}

func (d *testDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	d.recorder.WriteString(query)
	return sqlx.SelectContext(ctx, d.GetDatabase(), dest, query, args...)
}

func (d *testDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	d.recorder.WriteString(query)
	return d.GetDatabase().QueryxContext(ctx, query, args...)
}

func (d *testDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.recorder.WriteString(query)
	return d.GetDatabase().ExecContext(ctx, query, args...)
}

func sqlxRowsMapScan(rows *sqlx.Rows) []map[string]interface{} {
	res := []map[string]interface{}{}
	for rows.Next() {
		row := map[string]interface{}{}
		rows.MapScan(row)
		res = append(res, row)
	}

	return res
}
