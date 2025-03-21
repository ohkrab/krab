package krabdb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TransactionExecerContext interface {
	ExecerContext

	Rollback() error
	Commit() error
}

// Transaction represents database transaction.
type Transaction struct {
	tx *sqlx.Tx
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

func (t *Transaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// NullTransaction represents fake transaction.
type NullTransaction struct {
	db DB
}

func (t *NullTransaction) Rollback() error {
	return nil
}

func (t *NullTransaction) Commit() error {
	return nil
}

func (t *NullTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.db.GetDatabase().ExecContext(ctx, query, args...)
}
