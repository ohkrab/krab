package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	_ "embed"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/parser"
	"github.com/ohkrab/krab/ferro/run"
	"github.com/ohkrab/krab/ferro/run/generators"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/tpls"

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

	//go:embed tpls/embed/migration.fyml.tpl
	tplMigration []byte

	//go:embed tpls/embed/set.fyml.tpl
	tplSet []byte

	//go:embed tpls/embed/driver.fyml.tpl
	tplDriver []byte
)

func init() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s\n", cmd.Root().Version)
	}
}

func mustConfig(fs *config.Filesystem) *config.Config {
	parser := parser.New(fs)
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

	return cfg
}

func main() {
	templates := tpls.New(template.FuncMap{})
	templates.AddEmbedded("migration", tplMigration)
	templates.AddEmbedded("set", tplSet)
	templates.AddEmbedded("driver", tplDriver)

	dir, err := config.Dir()
	if err != nil {
		fmtx.WriteError("can't read config dir: %w", err)
		os.Exit(1)
	}

	filesystem := config.NewFilesystem(dir)

	// runners
	generator := run.NewGenerator(filesystem, templates, &generators.TimestampVersionGenerator{})
	migrator := run.NewMigrator(filesystem)

	// init commands
	initCmd := &cli.Command{
		Name:  "init",
		Usage: "Initialize default files structure",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return generator.GenInit(ctx, run.GenerateInitOptions{})
		},
	}

	// validate commands
	validateCmd := &cli.Command{
		Name:  "validate",
		Usage: "Validate the config",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			mustConfig(filesystem)
			fmtx.WriteSuccess("Config is valid")
			return nil
		},
	}

	// migrate commands
	migrateInitCmd := &cli.Command{
		Name:  "init",
		Usage: "Initialize a single timestamped migration",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return generator.GenMigration(ctx, run.GenerateMigrationOptions{})
		},
	}
	migrateAuditCmd := &cli.Command{
		Name:  "audit",
		Usage: "Show the audit logs",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := mustConfig(filesystem)
			return migrator.MigrateAudit(ctx, cfg, run.MigrateAuditOptions{})
		},
	}
	migrateUpCmd := &cli.Command{
		Name:  "up",
		Usage: "Apply all pending migrations",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := mustConfig(filesystem)
			return migrator.MigrateUp(ctx, cfg, run.MigrateUpOptions{})
		},
	}
	migrateDownCmd := &cli.Command{
		Name:  "down",
		Usage: "Rollback single migration",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := mustConfig(filesystem)
			return migrator.MigrateDown(ctx, cfg, run.MigrateDownOptions{})
		},
	}
	migrateStatusCmd := &cli.Command{
		Name:  "status",
		Usage: "Show the current status of migrations",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := mustConfig(filesystem)

			return migrator.MigrateStatus(ctx, cfg, run.MigrateStatusOptions{})
		},
	}
	migrateGroup := &cli.Command{
		Name:  "migrate",
		Usage: "Manage your migrations",
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
		validateCmd,
		migrateGroup,
	}

	err = root.Run(context.Background(), os.Args)
	if err != nil {
		fmtx.WriteError(err.Error())
		os.Exit(1)
	}
}
