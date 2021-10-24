package krab

import (
	"fmt"

	"github.com/ohkrab/krab/cli"
)

// ActionVersion prints full version.
type ActionVersion struct {
	Ui cli.UI
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
	a.Ui.Output(fmt.Sprint(InfoName, " ", InfoVersion))
	a.Ui.Output(fmt.Sprint("Build ", InfoCommit, " ", InfoBuildDate))
	return 0
}
