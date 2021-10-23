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

// ActionMigrateUp keeps data needed to perform this action.
type ActionMigrateUp struct {
	Set *MigrationSet
}

func (a *ActionMigrateUp) Help() string {
	return fmt.Sprint(
		`Usage: krab migrate up [set]`,
		"\n\n",
		a.Set.Arguments.Help(),
		`
Migrate all pending migrations in given [set].

Example:

    krab migrate up default
`,
	)
}

func (a *ActionMigrateUp) Synopsis() string {
	return fmt.Sprintf("Migrate `%s` up", a.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateUp) Run(args []string) int {
	ui := cli.DefaultUI()
	flags := cliargs.New(args)

	for _, arg := range a.Set.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	err = a.Set.Arguments.Validate(flags.Values())
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	templates := tpls.New(flags.Values())

	err = krabdb.WithConnection(func(db *sqlx.DB) error {
		return a.Do(context.Background(), db, templates, ui)
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
func (a *ActionMigrateUp) Do(ctx context.Context, db *sqlx.DB, tpl *tpls.Templates, ui cli.UI) error {
	versions := NewSchemaMigrationTable(tpl.Render(a.Set.Schema))

	// locking
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
	err = versions.Init(ctx, db)
	if err != nil {
		return errors.Wrap(err, "Failed to create default table for migrations")
	}

	migrationRefsInDb, err := versions.SelectAll(ctx, db)
	if err != nil {
		return err
	}

	pendingMigrations := versions.FilterPending(a.Set.Migrations, migrationRefsInDb)

	for _, pending := range pendingMigrations {
		ui.Output(fmt.Sprint(pending.RefName, " ", pending.Version))
		tx, err := krabdb.NewTx(ctx, db, pending.ShouldRunInTransaction())
		if err != nil {
			return errors.Wrap(err, "Failed to start transaction")
		}
		err = hooksRunner.SetSearchPath(ctx, tx, tpl.Render(a.Set.Schema))
		if err != nil {
			return errors.Wrap(err, "Failed to run SetSearchPath hook")
		}

		err = a.migrateUp(ctx, tx, pending, versions)
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

func (a *ActionMigrateUp) migrateUp(ctx context.Context, tx krabdb.TransactionExecerContext, migration *Migration, versions SchemaMigrationTable) error {
	_, err := tx.ExecContext(ctx, migration.Up.SQL)
	if err != nil {
		return errors.Wrap(err, "Failed to execute migration")
	}

	err = versions.Insert(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to insert migration")
	}

	return nil
}
