package krab

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
	"github.com/pkg/errors"
)

type ActionMigrateUp struct {
	Set *MigrationSet
	db  *sqlx.DB
}

func (a *ActionMigrateUp) migrate(ctx context.Context, migration *Migration) error {
	// BEGIN
	_, err := a.db.ExecContext(ctx, migration.Up.Sql)
	if err != nil {
		return err
	}

	err = SchemaMigrationInsert(ctx, a.db, migration.RefName)
	if err != nil {
		// ROLLBACK
		return err
	}

	// COMMIT
	return nil
}

func (a *ActionMigrateUp) Run(ctx context.Context) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	ok, err := krabdb.TryAdvisoryXactLock(ctx, tx, 1)
	if err != nil {
		return err
	}

	if ok {
		err := SchemaMigrationInit(ctx, a.db)
		if err != nil {
			return errors.Wrap(err, "Failed to create default table for migrations")
		}

		migrationRefsInDb, err := SchemaMigrationSelectAll(ctx, a.db)
		if err != nil {
			return err
		}

		pendingMigrations := SchemaMigrationFilterPending(a.Set.Migrations, migrationRefsInDb)

		{
			tx, err := a.db.BeginTxx(ctx, nil)
			if err != nil {
				return err
			}

			for _, pending := range pendingMigrations {
				err := a.migrate(ctx, pending)
				if err != nil {
					tx.Rollback() // ignore rollback error
					return err
				}
			}

			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("Another migration in progress")
	}

	return nil
}
