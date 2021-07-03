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

// ActionMigrateDown keeps data needed to perform this action.
type ActionMigrateDown struct {
	Set           *MigrationSet
	DownMigration SchemaMigration
}

func (a *ActionMigrateDown) Help() string {
	return `Usage: krab migrate down [set] [version]
  
Rollback migration in given [set] identified by [version].

Example:

    krab migrate down default 20060102150405
`
}

func (a *ActionMigrateDown) Synopsis() string {
	return fmt.Sprintf("Migrate `%s` down", a.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateDown) Run(args []string) int {
	ui := cli.DefaultUI()
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	err := flags.Parse(args)
	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	args = flags.Args()
	switch len(args) {
	case 1:
		a.DownMigration = SchemaMigration{args[0]}
	default:
		ui.Info(a.Help())
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

// Do performs the action.
// Schema migration must exist before running it.
func (a *ActionMigrateDown) Do(ctx context.Context, db *sqlx.DB) error {
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
