package krabdb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type QueryerContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type DB struct {
	db *sqlx.DB
}

func Connect(connectionString string) (*DB, error) {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlx.SelectContext(ctx, d.db, dest, query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.db.QueryxContext(ctx, query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}
