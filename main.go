package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabcli"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabenv"
)

func main() {
	ui := cli.DefaultUI()

	dir, err := krabenv.ConfigDif()
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

	conn := &krabdb.DefaultConnection{}

	registry := krab.CmdRegistry{}
	registry.Register(&krab.CmdVersion{})

	for _, set := range config.MigrationSets {
		localSet := set

		registry.Register(&krab.CmdMigrateStatus{
			Set:        localSet,
			Connection: conn,
		})
	}

	// agent := krabapi.Agent{Registry: registry}
	// agent.Run()

	c := krabcli.New(ui, os.Args[1:], config, &krabdb.DefaultConnection{})

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
