package krab

import (
	"context"
	"flag"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krabdb"
	"github.com/pkg/errors"
)

// ActionMigrateUp keeps data needed to perform this action.
type ActionMigrateUp struct {
	Set *MigrationSet
}

func (a *ActionMigrateUp) Help() string {
	return `Usage: krab migrate up [set]
  
Migrate all pending migrations in given [set].

Example:

    krab migrate up default
`
}

func (a *ActionMigrateUp) Synopsis() string {
	return fmt.Sprintf("Migrate `%s` up", a.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateUp) Run(args []string) int {
	ui := cli.DefaultUI()
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	err := flags.Parse(args)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	args = flags.Args()
	switch len(args) {
	case 0: // ok
	default:
		ui.Output(a.Help())
		ui.Error("Invalid number of arguments")
		return 1
	}

	err = krabdb.WithConnection(func(db *sqlx.DB) error {
		return a.Do(context.Background(), db, ui)
	})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Done")

	return 0
}

// Run performs the action. All pending migrations will be executed.
// Migration schema is created if does not exist.
func (a *ActionMigrateUp) Do(ctx context.Context, db *sqlx.DB, ui cli.UI) error {
	lockID := int64(1)

	_, err := krabdb.TryAdvisoryLock(ctx, db, lockID)
	if err != nil {
		return errors.Wrap(err, "Possibly another migration in progress")
	}
	defer krabdb.AdvisoryUnlock(ctx, db, lockID)

	err = SchemaMigrationInit(ctx, db)
	if err != nil {
		return errors.Wrap(err, "Failed to create default table for migrations")
	}

	migrationRefsInDb, err := SchemaMigrationSelectAll(ctx, db)
	if err != nil {
		return err
	}

	pendingMigrations := SchemaMigrationFilterPending(a.Set.Migrations, migrationRefsInDb)

	for _, pending := range pendingMigrations {
		ui.Output(fmt.Sprint(pending.RefName, " ", pending.Version))
		tx, err := krabdb.NewTx(ctx, db, pending.ShouldRunInTransaction())
		if err != nil {
			return errors.Wrap(err, "Failed to start transaction")
		}

		err = a.migrateUp(ctx, tx, pending)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ActionMigrateUp) migrateUp(ctx context.Context, tx krabdb.TransactionExecerContext, migration *Migration) error {
	_, err := tx.ExecContext(ctx, migration.Up.SQL)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = SchemaMigrationInsert(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to insert migration")
	}

	return nil
}
