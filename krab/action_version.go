package krab

import (
	"bytes"
	"context"
	"encoding/json"
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
	cmd := &CmdVersion{}

	buf := &bytes.Buffer{}
	err := cmd.Do(context.Background(), CmdOpts{Writer: buf})
	if err != nil {
		a.Ui.Error(err.Error())
		return 1
	}

	var response ResponseVersion
	json.Unmarshal(buf.Bytes(), &response)

	a.Ui.Output(response.Name)
	a.Ui.Output(response.Build)

	return 0
}
