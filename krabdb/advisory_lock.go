package krabdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TryAdvisoryXactLock will try to obtain transaction-level exclusive lock.
// It will not wait for the lock to become available. It will either obtain the lock immediately and return true,
// or return false if the lock cannot be acquired immediately.
// If acquired, is automatically released at the end of the current transaction and cannot be released explicitly.
//
// https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
//
func TryAdvisoryXactLock(ctx context.Context, tx *sqlx.Tx, id int64) (bool, error) {
	res, err := tx.QueryContext(ctx, "SELECT pg_try_advisory_xact_lock($1)", id)
	if err != nil {
		return false, errors.Wrap(err, "Failed to obtain advisory lock")
	}
	defer res.Close()

	for res.Next() {
		var success bool
		err = res.Scan(&success)
		return success, err
	}

	return false, errors.New("Failed to obtain advisory lock")
}

// TryAdvisoryLock will try to obtain session-level exclusive lock.
// This will either obtain the lock immediately and return true,
// or return false without waiting if the lock cannot be acquired immediately.
//
// https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
//
func TryAdvisoryLock(ctx context.Context, q sqlx.QueryerContext, id int64) (bool, error) {
	res, err := q.QueryContext(ctx, "SELECT pg_try_advisory_lock($1)", id)
	if err != nil {
		return false, errors.Wrap(err, "Failed to obtain advisory lock")
	}
	defer res.Close()

	for res.Next() {
		var success bool
		err = res.Scan(&success)
		return success, err
	}

	return false, errors.New("Failed to obtain advisory lock")
}

// TryAdvisoryLock will release a previously-acquired exclusive session-level advisory lock.
// Returns true if the lock is successfully released. If the lock was not held,
// false is returned, and in addition, an SQL warning will be reported by the server.
//
// https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
//
func AdvisoryUnlock(ctx context.Context, q sqlx.QueryerContext, id int64) (bool, error) {
	res, err := q.QueryContext(ctx, "SELECT pg_advisory_unlock($1)", id)
	if err != nil {
		return false, errors.Wrap(err, "Failed to release advisory lock")
	}
	defer res.Close()

	for res.Next() {
		var success bool
		err = res.Scan(&success)
		return success, err
	}

	return false, errors.New("Failed to release advisory lock")
}
