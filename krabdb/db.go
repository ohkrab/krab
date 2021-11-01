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
	NewTx(ctx context.Context, createTransaction bool) (TransactionExecerContext, error)
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

// NewTx is a helper that creates real transaction or null one based on createTransaction flag.
func (d *Instance) NewTx(ctx context.Context, createTransaction bool) (TransactionExecerContext, error) {
	if createTransaction {
		tx, err := d.GetDatabase().BeginTxx(ctx, nil)
		if err != nil {
			return nil, err
		}
		return &Transaction{tx: tx}, nil
	}

	return &NullTransaction{db: d}, nil
}
