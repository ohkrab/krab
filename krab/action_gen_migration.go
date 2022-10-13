package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
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
	resp, err := a.Cmd.Do(context.Background(), CmdOpts{Inputs{"args": args}})
	if err != nil {
		a.Ui.Error(err.Error())
		return 1
	}

	response := resp.(ResponseGenMigration)

	a.Ui.Output(response.Path)
	a.Ui.Output(response.Ref)

	return 0
}
