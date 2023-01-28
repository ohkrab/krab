package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
)

// ActionGenMigration generates migration file.
type ActionGenMigration struct {
	Ui  cli.UI
	Cmd *CmdGenMigration
}

func (a *ActionGenMigration) Help() string {
	return `Usage: krab gen migration

Generates migration file.
`
}

func (a *ActionGenMigration) Synopsis() string {
	return fmt.Sprintf("Generate migration file")
}

// Run in CLI.
func (a *ActionGenMigration) Run(args []string) int {
	ui := a.Ui
	flags := cliargs.New(args)

	for _, arg := range a.Cmd.Arguments().Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	resp, err := a.Cmd.Do(context.Background(), CmdOpts{Inputs: flags.Values()})
	if err != nil {
		a.Ui.Error(err.Error())
		return 1
	}

	response := resp.(ResponseGenMigration)

	a.Ui.Output("File generated:")
	a.Ui.Info(response.Path)
	a.Ui.Output("Don't forget to add your migration to migration_set:")
	a.Ui.Output(`
    migration_set "public" {
      migrations = [
        ...`)
	a.Ui.Info(fmt.Sprint("        ", response.Ref, ","))
	a.Ui.Output(`        ...
      ]
    }
`)

	return 0
}
