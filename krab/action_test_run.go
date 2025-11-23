package krab

import (
	"context"
	"fmt"

)

// ActionTestRun outputs test runner.
type ActionTestRun struct {
	Ui  cli.UI
	Cmd *CmdTestRun
}

func (a *ActionTestRun) Help() string {
	return fmt.Sprint(
		`Usage: krab test suite`,
		"\n\n",
		`
Starts a test suite.
`,
	)
}

func (a *ActionTestRun) Synopsis() string {
	return fmt.Sprintf("Test suite")
}

// Run in CLI.
func (a *ActionTestRun) Run(args []string) int {
	ui := a.Ui

	_, err := a.Cmd.Do(context.Background(), CmdOpts{})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Done")

	return 0
}
