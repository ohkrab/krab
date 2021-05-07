package krabdb

import "context"

type DbTx struct {
	Errs    chan<- error
	success bool
}

type DbTxCallback = func(tx *DbTx) error

func (db *DbTx) Begin(ctx context.Context) {
}

func (db *DbTx) Commit(ctx context.Context) {
	if !db.success {
		db.rollback(ctx)
	}
}

func (db *DbTx) rollback(ctx context.Context) {
}

func DbTransaction(ctx context.Context, fn DbTxCallback) {
	tx := DbTx{success: true}

	tx.Begin(ctx)
	defer tx.Commit(ctx)

	err := fn(&tx)
	if err != nil {
		tx.Errs <- err
		tx.success = false
	}
}
