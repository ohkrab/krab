package krab

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type ActionMigrateUp struct {
	Set *MigrationSet
}

func (a *ActionMigrateUp) migrate(ctx context.Context, tx pgx.Tx, migration *Migration) error {
	_, err := tx.Exec(ctx, migration.Up.Sql)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = SchemaMigrationInsert(ctx, tx.Conn(), migration.RefName)
	if err != nil {
		return errors.Wrap(err, "Failed to insert migration")
	}

	return nil
}

func (a *ActionMigrateUp) Run(ctx context.Context, conn *pgx.Conn) error {
	err := conn.BeginFunc(ctx, func(tx pgx.Tx) error {

		fmt.Println(conn.PgConn().IsBusy(), "busy")
		// ok, err := krabdb.TryAdvisoryXactLock(ctx, mainTx, 1)
		// if err != nil {
		// 	mainTx.Rollback()
		// 	return err
		// }

		// if ok {
		err := SchemaMigrationInit(ctx, tx.Conn())
		if err != nil {
			tx.Rollback(ctx)
			return errors.Wrap(err, "Failed to create default table for migrations")
		}

		migrationRefsInDb, err := SchemaMigrationSelectAll(ctx, tx.Conn())
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		pendingMigrations := SchemaMigrationFilterPending(a.Set.Migrations, migrationRefsInDb)

		for _, pending := range pendingMigrations {
			err := a.migrate(ctx, tx, pending)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		}

		// } else {
		// 	mainTx.Rollback()
		// 	return errors.New("Another migration in progress")
		// }

		// fmt.Println("commit")
		// mainTx.Commit()
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}
	return nil
}
