package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
	"github.com/pkg/errors"
)

// ActionMigrateDown keeps data needed to perform this action.
type ActionMigrateDown struct {
	Set           *MigrationSet
	DownMigration SchemaMigration
}

// Run performs the action.
// Schema migration must exist before running it.
func (a *ActionMigrateDown) Run(ctx context.Context, db *sqlx.DB) error {
	migration := a.Set.FindMigrationByVersion(a.DownMigration.Version)
	if migration == nil {
		return fmt.Errorf("Migration `%s` not found in `%s` set",
			a.DownMigration.Version,
			a.Set.RefName)
	}

	mainTx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}

	_, err = krabdb.TryAdvisoryXactLock(ctx, mainTx, 1)
	if err != nil {
		mainTx.Rollback()
		return errors.Wrap(err, "Possibly another migration in progress")
	}

	err = a.migrateDown(ctx, mainTx, migration)
	if err != nil {
		mainTx.Rollback()
		return err
	}

	err = mainTx.Commit()
	return err
}

func (a *ActionMigrateDown) migrateDown(ctx context.Context, tx *sqlx.Tx, migration *Migration) error {
	_, err := tx.ExecContext(ctx, migration.Down.SQL)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = SchemaMigrationDelete(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to delete migration")
	}

	return nil
}
