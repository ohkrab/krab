package spec

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/ferro/plugin"
)

type mockDBConnection struct {
	recorder         []string
	assertedSQLIndex int
	driver           plugin.Driver
}

func (m *mockDBConnection) Get(f func(db *testDB) error) error {
	conn, err := m.driver.Connect(context.Background(), nil)
	if err != nil {
		return err
	}
	defer m.driver.Disconnect(context.Background(), conn)

	return f(&testDB{recorder: &m.recorder, db: conn})
}

type testDB struct {
	db       plugin.DriverConnection
	recorder *[]string
}

func (d *testDB) GetDatabase() plugin.DriverConnection {
	return d.db
}

func (d *testDB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	*d.recorder = append(*d.recorder, query)
	// return sqlx.SelectContext(ctx, d.GetDatabase(), dest, query, args...)
	return nil
}

func (d *testDB) QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	*d.recorder = append(*d.recorder, query)
	// return d.GetDatabase().QueryxContext(ctx, query, args...)
	return nil, nil
}

func (d *testDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	*d.recorder = append(*d.recorder, query)
	// return d.GetDatabase().ExecContext(ctx, query, args...)
	return nil, nil
}

func (d *testDB) NewTx(ctx context.Context, createTransaction bool) (plugin.DriverQuery, error) {
	query := d.db.Query(plugin.DriverExecutionContext{})
	if createTransaction {
		tx, err := query.Begin(ctx)
		if err != nil {
			return nil, err
		}
		return &mockTransaction{tx: tx, recorder: d.recorder}, nil
	}

	return &mockNullTransaction{tx: query, recorder: d.recorder}, nil
}

func sqlxRowsMapScan(rows *sqlx.Rows) []map[string]any {
	res := []map[string]any{}
	for rows.Next() {
		row := map[string]any{}
		rows.MapScan(row)
		res = append(res, row)
	}

	return res
}

type mockTransaction struct {
	tx       plugin.DriverQuery
	recorder *[]string
}

func (t *mockTransaction) Begin(ctx context.Context) (plugin.DriverQuery, error) {
	return t.tx.Begin(ctx)
}

func (t *mockTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(context.Background())
}

func (t *mockTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(context.Background())
}

func (t *mockTransaction) Exec(ctx context.Context, query string, args ...any) error {
	*t.recorder = append(*t.recorder, query)
	return t.tx.Exec(ctx, query, args...)
}

type mockNullTransaction struct {
	tx       plugin.DriverQuery
	recorder *[]string
}

func (t *mockNullTransaction) Begin(ctx context.Context) (plugin.DriverQuery, error) {
	return t.tx, nil
}

func (t *mockNullTransaction) Rollback(ctx context.Context) error {
	return nil
}

func (t *mockNullTransaction) Commit(ctx context.Context) error {
	return nil
}

func (t *mockNullTransaction) Exec(ctx context.Context, query string, args ...any) error {
	*t.recorder = append(*t.recorder, query)
	return t.tx.Exec(ctx, query, args...)
}
