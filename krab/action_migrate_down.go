package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/wzshiming/ctc"
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

	cmd := &CmdMigrateDown{
		Set:        a.Set,
		Connection: a.Connection,
	}
	resp, err := cmd.Do(context.Background(), CmdOpts{Inputs: flags.Values()})
	result := resp.([]ResponseMigrateDown)

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	for _, status := range result {
		uiMigrationStatusFromResponseDown(ui, status)
	}

	return 0
}

func uiMigrationStatusFromResponseDown(ui cli.UI, resp ResponseMigrateDown) {
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
