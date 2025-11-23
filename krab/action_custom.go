package krab

import (
	"context"
	"fmt"

)

// ActionCustom keeps data needed to perform this action.
type ActionCustom struct {
	Ui  cli.UI
	Cmd *CmdAction
}

func (a *ActionCustom) Help() string {
	return fmt.Sprint(
		`Usage: krab action namespace name`,
		"\n\n",
		a.Cmd.Action.Arguments.Help(),
		`
Performs custom action.
`,
	)
}

func (a *ActionCustom) Synopsis() string {
	return fmt.Sprintf("Action")
}

// Run in CLI.
func (a *ActionCustom) Run(args []string) int {
	ui := a.Ui
	flags := cliargs.New(args)

	for _, arg := range a.Cmd.Action.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	_, err = a.Cmd.Do(context.Background(), CmdOpts{NamedInputs: flags.Values()})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Done")

	return 0
}
