package krabcli

import (
	"fmt"

	mcli "github.com/mitchellh/cli"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
)

type Command mcli.Command

type App struct {
	Ui     cli.UI
	CLI    *mcli.CLI
	Config *krab.Config
}

func New(ui cli.UI, args []string, config *krab.Config) *App {
	c := mcli.NewCLI(krab.InfoName, krab.InfoVersion)
	c.Args = args
	c.Commands = make(map[string]mcli.CommandFactory, 0)

	app := &App{
		Ui:     ui,
		CLI:    c,
		Config: config,
	}
	app.RegisterAll()

	return app
}

func (a *App) RegisterAll() {
	a.RegisterCmd("version", func() Command {
		return &krab.ActionVersion{Ui: a.Ui}
	})

	for _, set := range a.Config.MigrationSets {
		localSet := set

		a.RegisterCmd(fmt.Sprintln("migrate", "up", set.RefName), func() Command {
			return &krab.ActionMigrateUp{Ui: a.Ui, Set: localSet, Connection: &krabdb.DefaultConnection{}}
		})

		a.RegisterCmd(fmt.Sprintln("migrate", "down", set.RefName), func() Command {
			return &krab.ActionMigrateDown{
				Ui:         a.Ui,
				Set:        localSet,
				Connection: &krabdb.DefaultConnection{},
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
