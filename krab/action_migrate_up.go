package krab

import (
	"context"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
	"github.com/pkg/errors"
)

type ActionMigrateUp struct {
	Set *MigrationSet
}

func (a *ActionMigrateUp) migrate(ctx context.Context, tx *sqlx.Tx, migration *Migration) error {
	_, err := tx.ExecContext(ctx, migration.Up.Sql)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = SchemaMigrationInsert(ctx, tx, migration.RefName)
	if err != nil {
		return errors.Wrap(err, "Failed to insert migration")
	}

	return nil
}

func (a *ActionMigrateUp) Run(ctx context.Context, db *sqlx.DB) error {
	mainTx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}

	_, err = krabdb.TryAdvisoryXactLock(ctx, mainTx, 1)
	if err != nil {
		mainTx.Rollback()
		return errors.Wrap(err, "Possibly another migration in progress")
	}

	err = SchemaMigrationInit(ctx, mainTx)
	if err != nil {
		mainTx.Rollback()
		return errors.Wrap(err, "Failed to create default table for migrations")
	}

	migrationRefsInDb, err := SchemaMigrationSelectAll(ctx, mainTx)
	if err != nil {
		mainTx.Rollback()
		return err
	}

	pendingMigrations := SchemaMigrationFilterPending(a.Set.Migrations, migrationRefsInDb)

	for _, pending := range pendingMigrations {
		err := a.migrate(ctx, mainTx, pending)
		if err != nil {
			mainTx.Rollback()
			return err
		}
	}

	err = mainTx.Commit()
	return err
}
