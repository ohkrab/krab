package krabcli

import (
	"fmt"
	"strings"

	mcli "github.com/mitchellh/cli"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
)

type Command mcli.Command

type App struct {
	Ui         cli.UI
	CLI        *mcli.CLI
	Config     *krab.Config
	connection krabdb.Connection
}

func New(
	ui cli.UI,
	args []string,
	config *krab.Config,
	connection krabdb.Connection,
) *App {
	c := mcli.NewCLI(krab.InfoName, krab.InfoVersion)
	c.Args = args
	c.Commands = make(map[string]mcli.CommandFactory, 0)

	app := &App{
		Ui:         ui,
		CLI:        c,
		Config:     config,
		connection: connection,
	}
	app.RegisterAll()

	return app
}

func (a *App) RegisterAll() {
	a.RegisterCmd("version", func() Command {
		return &krab.ActionVersion{Ui: a.Ui}
	})

	for _, action := range a.Config.Actions {
		localAction := action
		a.RegisterCmd(strings.Join(action.Addr().Absolute(), " "), func() Command {
			return &krab.ActionCustom{Ui: a.Ui, Action: localAction, Connection: a.connection}
		})
	}

	for _, set := range a.Config.MigrationSets {
		localSet := set

		a.RegisterCmd(fmt.Sprintln("migrate", "status", set.RefName), func() Command {
			return &krab.ActionMigrateStatus{
				Ui:         a.Ui,
				Set:        localSet,
				Connection: a.connection,
			}
		})

		a.RegisterCmd(fmt.Sprintln("migrate", "up", set.RefName), func() Command {
			return &krab.ActionMigrateUp{Ui: a.Ui, Set: localSet, Connection: a.connection}
		})

		a.RegisterCmd(fmt.Sprintln("migrate", "down", set.RefName), func() Command {
			return &krab.ActionMigrateDown{
				Ui:         a.Ui,
				Set:        localSet,
				Connection: a.connection,
				Arguments: krab.Arguments{
					Args: []*krab.Argument{
						{
							Name:        "version",
							Type:        "string",
							Description: "Migration version to rollback",
						},
					},
				}}
		})
	}
}

func (a *App) Run() (int, error) {
	return a.CLI.Run()
}

func (a *App) RegisterCmd(names string, cmd func() Command) {
	a.CLI.Commands[names] = func() (mcli.Command, error) {
		return cmd(), nil
	}
}
