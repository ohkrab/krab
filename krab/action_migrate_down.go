package krab

import (
	"context"
	"fmt"

	"github.com/wzshiming/ctc"
)

// ActionMigrateDown keeps data needed to perform this action.
type ActionMigrateDown struct {
	Ui  cli.UI
	Cmd *CmdMigrateDown
}

func (a *ActionMigrateDown) Help() string {
	return fmt.Sprint(
		`Usage: krab migrate down [set] -version VERSION`,
		"\n\n",
		a.Cmd.Arguments().Help(),
		a.Cmd.Set.Arguments.Help(),
		`
Rollback migration in given [set] identified by VERSION.

Example:

    krab migrate down default -version 20060102150405
`,
	)
}

func (a *ActionMigrateDown) Synopsis() string {
	return fmt.Sprintf("Migrate `%s` down", a.Cmd.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateDown) Run(args []string) int {
	flags := cliargs.New(args)

	for _, arg := range a.Cmd.Set.Arguments.Args {
		flags.Add(arg.Name)
	}
	// default arguments always take precedence over custom ones
	for _, arg := range a.Cmd.Arguments().Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		a.Ui.Output(a.Help())
		a.Ui.Error(err.Error())
		return 1
	}

	resp, err := a.Cmd.Do(context.Background(), CmdOpts{NamedInputs: flags.Values()})
	result, ok := resp.([]ResponseMigrateDown)

	if err != nil {
		a.Ui.Error(err.Error())
		return 1
	}

	if ok {
		for _, status := range result {
			uiMigrationStatusFromResponseDown(a.Ui, status)
		}
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
