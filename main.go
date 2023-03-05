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

	dir, err := krabenv.ConfigDir()
	if err != nil {
		ui.Error(fmt.Errorf("can't read config dir: %w", err).Error())
		os.Exit(1)
	}

	parser := krab.NewParser()
	config, err := parser.LoadConfigDir(dir)
	if err != nil {
		ui.Error(fmt.Errorf("parsing error: %w", err).Error())
		os.Exit(1)
	}

	conn := &krabdb.DefaultConnection{}

	registry := &krab.CmdRegistry{
		Commands:         []krab.Cmd{},
		FS:               parser.FS,
		VersionGenerator: &krab.TimestampVersionGenerator{},
	}
	registry.RegisterAll(config, conn)

	c := krabcli.New(ui, os.Args[1:], config, registry, conn)

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
