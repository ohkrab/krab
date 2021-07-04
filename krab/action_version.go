package krab

import (
	"fmt"

	"github.com/ohkrab/krab/cli"
)

// ActionVersion prints full version.
type ActionVersion struct{}

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
	ui := cli.DefaultUI()
	ui.Output(fmt.Sprint(InfoName, " ", InfoVersion))
	ui.Output(fmt.Sprint("Build ", InfoCommit, " ", InfoBuildDate))
	return 0
}
