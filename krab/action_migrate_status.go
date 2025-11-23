package krab

import (
	"context"
	"fmt"
)

// ActionMigrateStatus keeps data needed to perform this action.
type ActionMigrateStatus struct {
	Ui  cli.UI
	Cmd *CmdMigrateStatus
}

func (a *ActionMigrateStatus) Help() string {
	return fmt.Sprint(
		`Usage: krab migrate status [set]`,
		"\n\n",
		a.Cmd.Set.Arguments.Help(),
		`
View migration status for given set.
`,
	)
}

func (a *ActionMigrateStatus) Synopsis() string {
	return fmt.Sprintf("Migration status for `%s`", a.Cmd.Set.RefName)
}

// Run in CLI.
func (a *ActionMigrateStatus) Run(args []string) int {
	ui := a.Ui
	flags := cliargs.New(args)

	for _, arg := range a.Cmd.Set.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	resp, err := a.Cmd.Do(context.Background(), CmdOpts{NamedInputs: flags.Values()})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	for _, status := range resp.([]ResponseMigrateStatus) {
		if status.Pending {
			ui.Output(cli.Red(fmt.Sprint("- ", status.Version, " ", status.Name)))
		} else {
			ui.Output(fmt.Sprint(emojis.CheckMark(), " ", status.Version, " ", status.Name))
		}
	}

	return 0
}
