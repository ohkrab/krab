package krab

import (
	"context"
	"fmt"

)

// ActionVersion prints full version.
type ActionVersion struct {
	Ui  cli.UI
	Cmd *CmdVersion
}

func (a *ActionVersion) Help() string {
	return `Usage: krab version

Prints full version.
`
}

func (a *ActionVersion) Synopsis() string {
	return fmt.Sprintf("Print full version")
}

// Run in CLI.
func (a *ActionVersion) Run(args []string) int {
	resp, err := a.Cmd.Do(context.Background(), CmdOpts{})
	if err != nil {
		a.Ui.Error(err.Error())
		return 1
	}

	response := resp.(ResponseVersion)

	a.Ui.Output(response.Name)
	a.Ui.Output(response.Build)

	return 0
}
