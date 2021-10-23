package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
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

	c := cli.New(krab.InfoName, krab.InfoVersion)

	c.RegisterCmd("version", func() cli.Command {
		return &krab.ActionVersion{}
	})

	for _, set := range config.MigrationSets {
		localSet := set

		c.RegisterCmd(fmt.Sprintln("migrate", "up", set.RefName), func() cli.Command {
			return &krab.ActionMigrateUp{Set: localSet}
		})

		c.RegisterCmd(fmt.Sprintln("migrate", "down", set.RefName), func() cli.Command {
			return &krab.ActionMigrateDown{Set: localSet, Arguments: krab.Arguments{
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

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
