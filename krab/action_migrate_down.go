package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/tpls"
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
	flags := cliargs.New(args)
	flags.RequireNonFlagArgs(1)

	for _, arg := range a.Set.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	args = flags.Args()
	switch len(args) {
	case 1:
		a.DownMigration = SchemaMigration{args[0]}
	default:
		err = krabdb.WithConnection(func(db *sqlx.DB) error {
			ui.Output("Latest migrations:")
			versions := NewSchemaMigrationTable(a.Set.Schema)
			migrations, err := versions.SelectLastN(context.TODO(), db, 5)
			for _, m := range migrations {
				ui.Info(fmt.Sprint("* ", m.Version))
			}
			ui.Output("")

			return err
		})
		if err != nil {
			ui.Error(err.Error())
			return 1
		}

		ui.Output(a.Help())
		ui.Error("Invalid number of arguments")
		return 1
	}

	templates := tpls.New(flags.Values())

	err = krabdb.WithConnection(func(db *sqlx.DB) error {
		return a.Do(context.Background(), db, templates)
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
func (a *ActionMigrateDown) Do(ctx context.Context, db *sqlx.DB, tpl *tpls.Templates) error {
	versions := NewSchemaMigrationTable(tpl.Render(a.Set.Schema))

	migration := a.Set.FindMigrationByVersion(a.DownMigration.Version)
	if migration == nil {
		return fmt.Errorf("Migration `%s` not found in `%s` set",
			a.DownMigration.Version,
			a.Set.RefName)
	}

	lockID := int64(1)

	_, err := krabdb.TryAdvisoryLock(ctx, db, lockID)
	if err != nil {
		return errors.Wrap(err, "Possibly another migration in progress")
	}
	defer krabdb.AdvisoryUnlock(ctx, db, lockID)

	hooksRunner := HookRunner{}
	err = hooksRunner.SetSearchPath(ctx, db, tpl.Render(a.Set.Schema))
	if err != nil {
		return errors.Wrap(err, "Failed to run SetSearchPath hook")
	}

	// schema migration
	tx, err := krabdb.NewTx(ctx, db, migration.ShouldRunInTransaction())
	if err != nil {
		return errors.Wrap(err, "Failed to start transaction")
	}
	err = hooksRunner.SetSearchPath(ctx, tx, tpl.Render(a.Set.Schema))
	if err != nil {
		return errors.Wrap(err, "Failed to run SetSearchPath hook")
	}

	migrationExists, _ := versions.Exists(ctx, db, SchemaMigration{migration.Version})
	if migrationExists {
		err = a.migrateDown(ctx, tx, migration, versions)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return errors.New("Migration has not been run yet, nothing to rollback")
	}

	err = tx.Commit()
	return err
}

func (a *ActionMigrateDown) migrateDown(ctx context.Context, tx krabdb.TransactionExecerContext, migration *Migration, versions SchemaMigrationTable) error {
	_, err := tx.ExecContext(ctx, migration.Down.SQL)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = versions.Delete(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to delete migration")
	}

	return nil
}
