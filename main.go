package main

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/parser"
	"github.com/ohkrab/krab/ferro/run"
	"github.com/ohkrab/krab/ferro/run/generators"
	"github.com/ohkrab/krab/fmtx"

	// "github.com/ohkrab/krab/cli"
	"github.com/urfave/cli/v3"
)

var (
	//go:embed res/favicon/favicon.ico
	favicon []byte

	//go:embed res/crab-final-pure-white.svg
	whiteLogo []byte

	//go:embed res/crab-final-pure.svg
	logo []byte
)

// func main() {
// 	ui := cli.DefaultUI()

// 	configDir, err := config.Dir()
// 	if err != nil {
// 		ui.Error(fmt.Errorf("can't read config dir: %w", err).Error())
// 		os.Exit(1)
// 	}

// 	preconfig, err := parser.New(configDir).LoadConfigDir(configDir)
// 	if err != nil {
// 		ui.Error(fmt.Errorf("parsing error: %w", err).Error())
// 		os.Exit(1)
// 	}
// 	cfg, err := preconfig.BuildConfig()
// 	if err != nil {
// 		ui.Error(fmt.Errorf("config: %w", err).Error())
// 		os.Exit(1)
// 	}

// 	for _, migrationSet := range cfg.MigrationSets {
// 		ui.Info(fmt.Sprintf("MigrationSet: %s", migrationSet.Metadata.Name))
// 	}

// 	// conn := &krabdb.DefaultConnection{}
// 	// switchableConn := &krabdb.SwitchableDatabaseConnection{}

// 	// srv := &web.Server{
// 	// 	Config:     config,
// 	// 	Connection: switchableConn,
// 	// 	EmbeddableResources: web.EmbeddableResources{
// 	// 		Favicon:   favicon,
// 	// 		WhiteLogo: whiteLogo,
// 	// 		Logo:      logo,
// 	// 	},
// 	// }

// 	// registry := &krab.CmdRegistry{
// 	// 	Commands:         []krab.Cmd{},
// 	// 	FS:               parser.FS,
// 	// 	VersionGenerator: &krab.TimestampVersionGenerator{},
// 	// }
// 	// registry.RegisterAll(config, conn)

// 	// c := krabcli.New(ui, os.Args[1:], config, registry, conn, srv)

// 	// exitStatus, err := c.Run()
// 	// if err != nil {
// 	// 	ui.Error(err.Error())
// 	// }

// 	// os.Exit(exitStatus)
// }

func init() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s\n", cmd.Root().Version)
	}
}

func main() {
	dir, err := config.Dir()
	if err != nil {
		fmtx.WriteError("can't read config dir: %w", err)
		os.Exit(1)
	}

	filesystem := config.NewFilesystem(dir)
	parser := parser.New(filesystem)
	parsed, err := parser.LoadAndParse()
	if err != nil {
		fmtx.WriteError(err.Error())
		os.Exit(1)
	}

	cfg, errs := parsed.BuildConfig()
	if errs.HasErrors() {
		for _, err := range errs.Errors {
			fmtx.WriteError(err.Error())
		}
		os.Exit(1)
	}

	// runners
	initializer := run.NewInitializer(filesystem, &generators.TimestampVersionGenerator{})
	migrator := run.NewMigrator(filesystem)

	// init commands
	initCmd := &cli.Command{
		Name:  "init",
		Usage: "Initialize FerroDB default files structure",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return initializer.Initialize(ctx, run.InitializeOptions{})
		},
	}

	// migrate commands
	migrateInitCmd := &cli.Command{
		Name:  "init",
		Usage: "Initialize a single timestamped migration",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return migrator.MigrateInit(ctx, cfg, run.MigrateInitOptions{})
		},
	}
	migrateAuditCmd := &cli.Command{
		Name:  "audit",
		Usage: "Show the audit logs",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return migrator.MigrateAudit(ctx, cfg, run.MigrateAuditOptions{})
		},
	}
	migrateUpCmd := &cli.Command{
		Name:  "up",
		Usage: "Apply all pending migrations",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return migrator.MigrateUp(ctx, cfg, run.MigrateUpOptions{})
		},
	}
	migrateDownCmd := &cli.Command{
		Name:  "down",
		Usage: "Rollback single migration",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return migrator.MigrateDown(ctx, cfg, run.MigrateDownOptions{})
		},
	}
	migrateStatusCmd := &cli.Command{
		Name:  "status",
		Usage: "Show the current status of migrations",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return migrator.MigrateStatus(ctx, cfg, run.MigrateStatusOptions{})
		},
	}
	migrateGroup := &cli.Command{
		Name: "migrate",
		Commands: []*cli.Command{
			migrateInitCmd,
			migrateAuditCmd,
			migrateStatusCmd,
			migrateUpCmd,
			migrateDownCmd,
		},
	}

	root := &cli.Command{
		Name:      "ferro",
		Version:   "v0.20.0",
		Copyright: "(c) @qbart",
		Usage:     "ferroDB is a tool for managing your databases",
		UsageText: "ferro [global options] command [command options] [arguments...]",
	}
	root.Commands = []*cli.Command{
		initCmd,
		migrateGroup,
	}

	err = root.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(1)
	}
}
