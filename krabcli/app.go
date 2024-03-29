package krabcli

import (
	"fmt"
	"strings"

	mcli "github.com/mitchellh/cli"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/web"
)

type Command mcli.Command

type App struct {
	Ui         cli.UI
	CLI        *mcli.CLI
	Config     *krab.Config
	Registry   *krab.CmdRegistry
	connection krabdb.Connection
}

func New(
	ui cli.UI,
	args []string,
	config *krab.Config,
	registry *krab.CmdRegistry,
	connection krabdb.Connection,
	srv *web.Server,
) *App {
	c := mcli.NewCLI(krab.InfoName, krab.InfoVersion)
	c.Args = args
	c.Commands = make(map[string]mcli.CommandFactory, 0)

	app := &App{
		Ui:         ui,
		CLI:        c,
		Config:     config,
		Registry:   registry,
		connection: connection,
	}
	app.RegisterAll(srv)

	return app
}

func (a *App) RegisterAll(srv *web.Server) {
	a.RegisterCmd("server", func() Command {
		return srv
	})

	for _, cmd := range a.Registry.Commands {
		name := strings.Join(cmd.Name(), " ")

		switch c := cmd.(type) {
		case *krab.CmdVersion:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionVersion{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdMigrateDown:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionMigrateDown{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdMigrateUp:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionMigrateUp{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdMigrateStatus:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionMigrateStatus{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdAction:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionCustom{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdTestRun:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionTestRun{Ui: a.Ui, Cmd: c}
			})

		case *krab.CmdGenMigration:
			a.RegisterCmd(name, func() Command {
				return &krab.ActionGenMigration{Ui: a.Ui, Cmd: c}
			})

		default:
			panic(fmt.Sprintf("Not implemented: failed to register CLI action for command %T", c))
		}
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
