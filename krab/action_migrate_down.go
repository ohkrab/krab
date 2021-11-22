package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/tpls"
	"github.com/pkg/errors"
)

// ActionMigrateDown keeps data needed to perform this action.
type ActionMigrateDown struct {
	Ui            cli.UI
	Set           *MigrationSet
	DownMigration SchemaMigration
	Arguments     Arguments
	Connection    krabdb.Connection
}

func (a *ActionMigrateDown) Help() string {
	return fmt.Sprint(
		`Usage: krab migrate down [set] -version VERSION`,
		"\n\n",
		a.Arguments.Help(),
		a.Set.Arguments.Help(),
		` 
Rollback migration in given [set] identified by VERSION.

Example:

    krab migrate down default -version 20060102150405
`,
	)
}

func (a *ActionMigrateDown) Synopsis() string {
	return fmt.Sprintf("Migrate `%s` down", a.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateDown) Run(args []string) int {
	ui := a.Ui
	flags := cliargs.New(args)
	flags.RequireNonFlagArgs(0)

	for _, arg := range a.Set.Arguments.Args {
		flags.Add(arg.Name)
	}
	// default arguments always take precedence over custom ones
	for _, arg := range a.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}
	err = a.Arguments.Validate(flags.Values())
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

	a.DownMigration = SchemaMigration{
		cliargs.Values(flags.Values()).Get("version"),
	}

	// err = krabdb.WithConnection(func(db *sqlx.DB) error {
	// 	ui.Output("Latest migrations:")
	// 	versions := NewSchemaMigrationTable(a.Set.Schema)
	// 	migrations, err := versions.SelectLastN(context.TODO(), db, 5)
	// 	for _, m := range migrations {
	// 		ui.Info(fmt.Sprint("* ", m.Version))
	// 	}
	// 	ui.Output("")

	// 	return err
	// })
	// if err != nil {
	// 	ui.Error(err.Error())
	// 	return 1
	// }

	// ui.Output(a.Help())
	// ui.Error("Invalid number of arguments")
	// return 1

	err = a.Connection.Get(func(db krabdb.DB) error {
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
func (a *ActionMigrateDown) Do(ctx context.Context, db krabdb.DB, tpl *tpls.Templates) error {
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
	tx, err := db.NewTx(ctx, migration.ShouldRunInTransaction())
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
	sqls := migration.Down.ToSQLStatements()
	for _, sql := range sqls {
		// fmt.Println(ctc.ForegroundYellow, string(sql), ctc.Reset)
		_, err := tx.ExecContext(ctx, string(sql))
		if err != nil {
			return errors.Wrap(err, "Failed to execute migration")
		}
	}

	err := versions.Delete(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to delete migration")
	}

	return nil
}
