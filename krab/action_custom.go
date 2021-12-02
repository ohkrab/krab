package krab

import (
	"context"
	"fmt"
	"strings"

	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
)

// ActionCustom keeps data needed to perform this action.
type ActionCustom struct {
	Ui         cli.UI
	Action     *Action
	Connection krabdb.Connection
}

func (a *ActionCustom) Help() string {
	return fmt.Sprint(
		`Usage: krab action namespace name`,
		"\n\n",
		a.Action.Arguments.Help(),
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
	flags.RequireNonFlagArgs(0)

	for _, arg := range a.Action.Arguments.Args {
		flags.Add(arg.Name)
	}

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	err = a.Action.Arguments.Validate(flags.Values())
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}

	templates := tpls.New(flags.Values(), krabtpl.Functions)

	err = a.Connection.Get(func(db krabdb.DB) error {
		return a.Do(context.Background(), db, templates)
	})

	if err != nil {
		ui.Error(err.Error())
		return 1
	}

	ui.Info("Done")

	return 0
}

// Do performs the action.
func (a *ActionCustom) Do(ctx context.Context, db krabdb.DB, tpl *tpls.Templates) error {
	sb := strings.Builder{}
	a.Action.ToSQL(&sb)
	sql := tpl.Render(sb.String())

	_, err := db.ExecContext(ctx, sql)

	return err
}
