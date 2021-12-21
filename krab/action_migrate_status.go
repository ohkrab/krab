package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/emojis"
	"github.com/ohkrab/krab/krabdb"
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

	for _, arg := range a.Set.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	cmd := &CmdMigrateStatus{
		Set:        a.Set,
		Connection: a.Connection,
	}
	resp, err := cmd.Do(context.Background(), CmdOpts{Inputs: flags.Values()})

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
