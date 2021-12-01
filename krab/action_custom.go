package krab

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/tpls"
)

// ActionCustom keeps data needed to perform this action.
type ActionCustom struct {
	Ui         cli.UI
	Action     *Action
	Arguments  Arguments
	Connection krabdb.Connection
}

func (a *ActionCustom) Help() string {
	return fmt.Sprint(
		`Usage: krab action namespace name`,
		"\n\n",
		a.Arguments.Help(),
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
	return 0
}

// Do performs the action.
func (a *ActionCustom) Do(ctx context.Context, db krabdb.DB, tpl *tpls.Templates) error {
	return nil
}
