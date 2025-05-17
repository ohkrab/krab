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
	"github.com/ohkrab/krab/ferro/run"
	"github.com/ohkrab/krab/ferro/run/generators"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/plugins"
	"github.com/ohkrab/krab/tpls"

	"github.com/urfave/cli/v3"
)

var (
	//go:embed res/ferrodbicon.svg
	favicon []byte

	//go:embed res/ferrodb.svg
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

func main() {
	// check config dir
	dir, err := config.Dir()
	if err != nil {
		fmtx.WriteError("can't read config dir: %w", err)
		os.Exit(1)
	}

	// init templates
	templates := tpls.New(template.FuncMap{})
	templates.AddEmbedded("migration", tplMigration)
	templates.AddEmbedded("set", tplSet)
	templates.AddEmbedded("driver", tplDriver)

	// init plugins
	registry := plugins.New()
	registry.RegisterAll()

	// init internals
	filesystem := config.NewFilesystem(dir)
	// runners
	generator := run.NewGenerator(filesystem, templates, &generators.TimestampVersionGenerator{})
	runner := run.New(filesystem, templates, registry)

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
			runner.MustConfig()
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
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
			&cli.UintFlag{Name: "n", Usage: "Last N events to show", Required: false, Aliases: []string{"n"}, DefaultText: "0"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			run := run.CommandMigrateAudit{
				Driver: cmd.String("driver"),
				Set:    cmd.String("set"),
				N:      uint(cmd.Uint("n")),
			}
            audited, err := runner.ExecuteMigrateAudit(ctx, &run)
            if err != nil {
                return err
            }
            fmtx.WriteSuccess("Audit logs for %s/%s", run.Driver, run.Set)
            for name, set := range audited.Sets {
                fmtx.WriteSuccess("Set: %s, Migrations: %d", name, len(set.Migrations))
            }

            return nil
		},
	}
	migrateUpCmd := &cli.Command{
		Name:  "up",
		Usage: "Apply all pending migrations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			run := run.CommandMigrateUp{
				Driver: cmd.String("driver"),
				Set:    cmd.String("set"),
			}
			return runner.ExecuteMigrateUp(ctx, &run)
		},
	}
	migrateDownCmd := &cli.Command{
		Name:  "down",
		Usage: "Rollback single migration",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
			&cli.StringFlag{Name: "version", Usage: "Version to rollback to", Required: true, Aliases: []string{"v"}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			run := run.CommandMigrateDown{
				Driver:  cmd.String("driver"),
				Set:     cmd.String("set"),
				Version: cmd.String("version"),
			}
			return runner.ExecuteMigrateDown(ctx, &run)
		},
	}
	migrateStatusCmd := &cli.Command{
		Name:  "status",
		Usage: "Show the current status of migrations",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			run := run.CommandMigrateStatus{
				Driver: cmd.String("driver"),
				Set:    cmd.String("set"),
			}
			return runner.ExecuteMigrateStatus(ctx, &run)
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
