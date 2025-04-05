package main

import (
	"fmt"
	"os"

	_ "embed"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/parser"
)

var (
	//go:embed res/favicon/favicon.ico
	favicon []byte

	//go:embed res/crab-final-pure-white.svg
	whiteLogo []byte

	//go:embed res/crab-final-pure.svg
	logo []byte
)

func main() {
	ui := cli.DefaultUI()

	configDir, err := config.Dir()
	if err != nil {
		ui.Error(fmt.Errorf("can't read config dir: %w", err).Error())
		os.Exit(1)
	}

	preconfig, err := parser.New(configDir).LoadConfigDir(configDir)
	if err != nil {
		ui.Error(fmt.Errorf("parsing error: %w", err).Error())
		os.Exit(1)
	}
	cfg, err := preconfig.BuildConfig()
	if err != nil {
		ui.Error(fmt.Errorf("config: %w", err).Error())
		os.Exit(1)
	}

	for _, migrationSet := range cfg.MigrationSets {
		ui.Info(fmt.Sprintf("MigrationSet: %s", migrationSet.Metadata.Name))
	}

	// conn := &krabdb.DefaultConnection{}
	// switchableConn := &krabdb.SwitchableDatabaseConnection{}

	// srv := &web.Server{
	// 	Config:     config,
	// 	Connection: switchableConn,
	// 	EmbeddableResources: web.EmbeddableResources{
	// 		Favicon:   favicon,
	// 		WhiteLogo: whiteLogo,
	// 		Logo:      logo,
	// 	},
	// }

	// registry := &krab.CmdRegistry{
	// 	Commands:         []krab.Cmd{},
	// 	FS:               parser.FS,
	// 	VersionGenerator: &krab.TimestampVersionGenerator{},
	// }
	// registry.RegisterAll(config, conn)

	// c := krabcli.New(ui, os.Args[1:], config, registry, conn, srv)

	// exitStatus, err := c.Run()
	// if err != nil {
	// 	ui.Error(err.Error())
	// }

	// os.Exit(exitStatus)
}
