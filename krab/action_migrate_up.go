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
		return a.Do(context.Background(), db)
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
func (a *ActionMigrateUp) Do(ctx context.Context, db *sqlx.DB) error {
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
		err := a.migrateUp(ctx, mainTx, pending)
		if err != nil {
			mainTx.Rollback()
			return err
		}
	}

	err = mainTx.Commit()
	return err
}

func (a *ActionMigrateUp) migrateUp(ctx context.Context, tx *sqlx.Tx, migration *Migration) error {
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
