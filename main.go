package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabcli"
	"github.com/ohkrab/krab/krabenv"
)

func main() {
	ui := cli.DefaultUI()

	dir, err := krabenv.GetConfigDir()
	if err != nil {
		ui.Error(fmt.Errorf("Can't read config dir: %w", err).Error())
		os.Exit(1)
	}

	parser := krab.NewParser()
	config, err := parser.LoadConfigDir(dir)
	if err != nil {
		ui.Error(fmt.Errorf("Parsing error: %w", err).Error())
		os.Exit(1)
	}

	c := krabcli.New(ui, os.Args[1:], config)

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
