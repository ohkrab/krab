package cli

import (
	"os"

	mcli "github.com/mitchellh/cli"
)

type Command mcli.Command

type App struct {
	cli *mcli.CLI
}

func New(name, version string) *App {
	c := mcli.NewCLI(name, version)
	c.Args = os.Args[1:]
	c.Commands = make(map[string]mcli.CommandFactory, 0)

	app := &App{
		cli: c,
	}
	return app
}

func (a *App) Run() (int, error) {
	return a.cli.Run()
}

func (a *App) RegisterCmd(names string, cmd func() Command) {
	a.cli.Commands[names] = func() (mcli.Command, error) {
		return cmd(), nil
	}
}
