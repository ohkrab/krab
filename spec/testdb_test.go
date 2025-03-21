package spec

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabenv"
)

type mockDBConnection struct {
	recorder         []string
	assertedSQLIndex int
}

func (m *mockDBConnection) Get(f func(db krabdb.DB) error) error {
	db, err := krabdb.Connect(krabenv.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close()

	return f(&testDB{recorder: &m.recorder, db: db})
}

type testDB struct {
	db       *sqlx.DB
	recorder *[]string
}

func (d *testDB) GetDatabase() *sqlx.DB {
	return d.db
}

func (d *testDB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	*d.recorder = append(*d.recorder, query)
	return sqlx.SelectContext(ctx, d.GetDatabase(), dest, query, args...)
}

func (d *testDB) QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	*d.recorder = append(*d.recorder, query)
	return d.GetDatabase().QueryxContext(ctx, query, args...)
}

func (d *testDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	*d.recorder = append(*d.recorder, query)
	return d.GetDatabase().ExecContext(ctx, query, args...)
}

func (d *testDB) NewTx(ctx context.Context, createTransaction bool) (krabdb.TransactionExecerContext, error) {
	if createTransaction {
		tx, err := d.GetDatabase().BeginTxx(ctx, nil)
		if err != nil {
			return nil, err
		}
		return &mockTransaction{tx: tx, recorder: d.recorder}, nil
	}

	return &mockNullTransaction{db: d, recorder: d.recorder}, nil
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
	tx       *sqlx.Tx
	recorder *[]string
}

func (t *mockTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *mockTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *mockTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	*t.recorder = append(*t.recorder, query)
	return t.tx.ExecContext(ctx, query, args...)
}

type mockNullTransaction struct {
	db       krabdb.DB
	recorder *[]string
}

func (t *mockNullTransaction) Rollback() error {
	return nil
}

func (t *mockNullTransaction) Commit() error {
	return nil
}

func (t *mockNullTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	*t.recorder = append(*t.recorder, query)
	return t.db.ExecContext(ctx, query, args...)
}
