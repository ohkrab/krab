package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/wzshiming/ctc"
)

// ActionMigrateUp keeps data needed to perform this action.
type ActionMigrateUp struct {
	Ui         cli.UI
	Set        *MigrationSet
	Connection krabdb.Connection
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
	ui := a.Ui
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

	cmd := &CmdMigrateUp{
		Set:        a.Set,
		Connection: a.Connection,
		Inputs:     flags.Values(),
	}
	resp, err := cmd.Do(context.Background(), CmdOpts{})
	result := resp.([]ResponseMigrateUp)

	if len(result) > 0 {
		for _, status := range result {
			uiMigrationStatusFromResponseUp(ui, status)
		}
	}

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	if len(result) == 0 {
		ui.Info("No pending migrations")
	}

	return 0
}

func uiMigrationStatusFromResponseUp(ui cli.UI, resp ResponseMigrateUp) {
	color := ctc.ForegroundGreen
	text := "OK  "
	if !resp.Success {
		color = ctc.ForegroundRed
		text = "ERR "
	}

	ui.Output(fmt.Sprint(
		color,
		text,
		ctc.Reset,
		resp.Version,
		" ",
		resp.Name,
	))
}
