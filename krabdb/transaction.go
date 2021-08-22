package krabdb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TransactionExecerContext interface {
	sqlx.ExecerContext

	Rollback() error
	Commit() error
}

// NewTx is a helper that creates real transaction or null one based on createTransaction flag.
func NewTx(ctx context.Context, db *sqlx.DB, createTransaction bool) (TransactionExecerContext, error) {
	if createTransaction {
		return BeginTx(ctx, db)
	}

	return NullTx(ctx, db)
}

// BeginTx starts new transaction
func BeginTx(ctx context.Context, db *sqlx.DB) (TransactionExecerContext, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
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

func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// NullTx returns fake transaction to satisfy TransactionExecerContext interface
func NullTx(ctx context.Context, db *sqlx.DB) (TransactionExecerContext, error) {
	return &NullTransaction{db: db}, nil
}

// NullTransaction represents fake transaction.
type NullTransaction struct {
	db *sqlx.DB
}

func (t *NullTransaction) Rollback() error {
	return nil
}

func (t *NullTransaction) Commit() error {
	return nil
}

func (t *NullTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.db.ExecContext(ctx, query, args...)
}
