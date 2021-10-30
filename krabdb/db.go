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

type DB interface {
	ExecerContext
	QueryerContext

	GetDatabase() *sqlx.DB
}

type Instance struct {
	DB
	database *sqlx.DB
}

func (d *Instance) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlx.SelectContext(ctx, d.GetDatabase(), dest, query, args...)
}

func (d *Instance) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.GetDatabase().QueryxContext(ctx, query, args...)
}

func (d *Instance) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.GetDatabase().ExecContext(ctx, query, args...)
}

func (d *Instance) GetDatabase() *sqlx.DB {
	return d.database
}
