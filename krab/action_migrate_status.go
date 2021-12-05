package krab

import (
	"context"
	"fmt"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
	"github.com/pkg/errors"
	"github.com/wzshiming/ctc"
)

// ActionMigrateStatus keeps data needed to perform this action.
type ActionMigrateStatus struct {
	Ui         cli.UI
	Set        *MigrationSet
	Connection krabdb.Connection
}

func (a *ActionMigrateStatus) Help() string {
	return fmt.Sprint(
		`Usage: krab migrate status [set]`,
		"\n\n",
		a.Set.Arguments.Help(),
		`
View migration status for given set.
`,
	)
}

func (a *ActionMigrateStatus) Synopsis() string {
	return fmt.Sprintf("Migration status for `%s`", a.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateStatus) Run(args []string) int {
	ui := a.Ui
	flags := cliargs.New(args)
	flags.RequireNonFlagArgs(0)

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

	templates := tpls.New(flags.Values(), krabtpl.Functions)

	err = a.Connection.Get(func(db krabdb.DB) error {
		return a.Do(context.Background(), db, templates, ui)
	})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	return 0
}

// Run performs the action.
func (a *ActionMigrateStatus) Do(ctx context.Context, db krabdb.DB, tpl *tpls.Templates, ui cli.UI) error {
	versions := NewSchemaMigrationTable(tpl.Render(a.Set.Schema))

	hooksRunner := HookRunner{}
	err := hooksRunner.SetSearchPath(ctx, db, tpl.Render(a.Set.Schema))
	if err != nil {
		return errors.Wrap(err, "Failed to run SetSearchPath hook")
	}
	migrationRefsInDb, err := versions.SelectAll(ctx, db)
	if err != nil {
		return err
	}

	appliedMigrations := hashset.New()

	for _, migration := range migrationRefsInDb {
		appliedMigrations.Add(migration.Version)
	}

	for _, migration := range a.Set.Migrations {
		pending := !appliedMigrations.Contains(migration.Version)

		if pending {
			ui.Error(fmt.Sprint("- ", migration.Version, " ", migration.RefName))
		} else {
			ui.Output(fmt.Sprint(ctc.ForegroundGreen, "+ ", ctc.Reset, migration.Version, " ", migration.RefName))
		}

	}

	return nil
}
