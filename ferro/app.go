package ferro

import (
	"context"
	"fmt"
	"text/template"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/run"
	"github.com/ohkrab/krab/ferro/run/generators"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/plugins"
	"github.com/ohkrab/krab/tpls"
	"github.com/urfave/cli/v3"
)

func init() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s\n", cmd.Root().Version)
	}
}

type App struct {
	EmbededDriverTemplate    []byte
	EmbededSetTemplate       []byte
	EmbededMigrationTemplate []byte

	Dir    string
	Logger *fmtx.Logger
}

func (a *App) Run(args []string) int {
	// init templates
	templates := tpls.New(template.FuncMap{})
	templates.AddEmbedded("migration", a.EmbededMigrationTemplate)
	templates.AddEmbedded("set", a.EmbededSetTemplate)
	templates.AddEmbedded("driver", a.EmbededDriverTemplate)

	// init plugins
	registry := plugins.New()
	registry.RegisterAll()

	// init internals
	filesystem := config.NewFilesystem(a.Dir)
	// runners
	generator := run.NewGenerator(filesystem, templates, &generators.TimestampVersionGenerator{})
	runner := run.New(filesystem, templates, registry, a.Logger)

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
            _, err := runner.UseConfig()
            if err != nil {
                return err
            }
			a.Logger.WriteSuccess("Config is valid")
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

	fixMigrationUp := &cli.Command{
		Name:  "up",
		Usage: "Mark UP migration as fixed",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
			&cli.StringFlag{Name: "version", Usage: "Version to fix", Required: true, Aliases: []string{"v"}},
			&cli.StringFlag{Name: "comment", Usage: "Comment for the fix", Required: false, Aliases: []string{"C"}, DefaultText: "Manually fixed", Value: "Manually fixed"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			runCmd := run.CommandMigrateFixUp{
				Driver:  cmd.String("driver"),
				Set:     cmd.String("set"),
				Version: cmd.String("version"),
				Comment: cmd.String("comment"),
			}
			_, err := runner.ExecuteMigrateFixUp(ctx, &runCmd)
			if err != nil {
				return err
			}
			a.Logger.WriteSuccess("Marked as fixed successfully")

			return nil
		},
	}
	fixMigrationDown := &cli.Command{
		Name:  "down",
		Usage: "Mark DOWN migration as fixed",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
			&cli.StringFlag{Name: "version", Usage: "Version to fix", Required: true, Aliases: []string{"v"}},
			&cli.StringFlag{Name: "comment", Usage: "Comment for the fix", Required: false, Aliases: []string{"C"}, DefaultText: "Manually fixed", Value: "Manually fixed"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			runCmd := run.CommandMigrateFixDown{
				Driver:  cmd.String("driver"),
				Set:     cmd.String("set"),
				Version: cmd.String("version"),
				Comment: cmd.String("comment"),
			}
			_, err := runner.ExecuteMigrateFixDown(ctx, &runCmd)
			if err != nil {
				return err
			}
			a.Logger.WriteSuccess("Marked as fixed successfully")

			return nil
		},
	}
	migrateFixGroup := &cli.Command{
		Name:  "fix",
		Usage: "Apply fixes to migrations",
		Commands: []*cli.Command{
			fixMigrationUp,
			fixMigrationDown,
		},
	}

	migrateAuditCmd := &cli.Command{
		Name:  "audit",
		Usage: "Show the audit logs",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "driver", Usage: "The driver to use", Required: true, Aliases: []string{"d"}},
			&cli.StringFlag{Name: "set", Usage: "MigrationSet to use", Required: true, Aliases: []string{"s"}},
			&cli.UintFlag{Name: "n", Usage: "Last N events to show", Required: false, Aliases: []string{"n"}, DefaultText: "0", Value: 0},
			&cli.StringFlag{Name: "f", Usage: "Format", Required: false, Aliases: []string{"f"}, DefaultText: "short", Value: "short", Validator: func(s string) error {
				if s != "short" && s != "long" {
					return fmt.Errorf("invalid format: %s, must be 'short' or 'long'", s)
				}
				return nil
			}},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			runCmd := run.CommandMigrateAudit{
				Driver:   cmd.String("driver"),
				Set:      cmd.String("set"),
				N:        uint(cmd.Uint("n")),
				FullView: cmd.String("f") == "long",
			}
			out, err := runner.ExecuteMigrateAudit(ctx, &runCmd)
			if err != nil {
				return err
			}
			if len(out.Logs) == 0 {
				a.Logger.WriteSuccess("No audit logs so far")
			} else {
				for _, log := range out.Logs {
					event := log.Event
					switch log.Event {
					case run.MigrationUpCompletedEvent, run.MigrationDownCompletedEvent:
						event = fmtx.Success("%-24s", log.Event)

					case run.MigrationUpFailedEvent, run.MigrationDownFailedEvent:
						event = fmtx.Danger("%-24s", log.Event)

					case run.MigrationFixUpEvent, run.MigrationFixDownEvent:
						event = fmtx.Warning("%-24s", log.Event)

					default:
						event = fmt.Sprintf("%-24s", log.Event)
					}

					sign := "  "
					if log.Event == run.MigrationUpCompletedEvent {
						sign = fmtx.Success("%-2s", "+")
					} else if log.Event == run.MigrationDownCompletedEvent {
						sign = fmtx.Danger("%-2s", "-")
					} else if log.Event == run.MigrationUpFailedEvent || log.Event == run.MigrationDownFailedEvent {
						sign = fmtx.Danger("%-2s", "!")
					} else if log.Event == run.MigrationFixUpEvent || log.Event == run.MigrationFixDownEvent {
						sign = fmtx.Warning("%-2s", "~")
					}

					a.Logger.WriteLine(
						"%s %3d %s %s %s %s",
						sign,
						log.ID,
						log.AppliedAt.Format("2006-01-02 15:04:05"),
						event,
						log.GetData("migration"),
						log.GetData("version"),
					)
					if runCmd.FullView {
						switch log.Event {
						case run.MigrationUpFailedEvent, run.MigrationDownFailedEvent:
							a.Logger.WriteLine(
								"    %s",
								fmtx.Danger("%s", log.GetMetadata("error")),
							)
						case run.MigrationFixUpEvent, run.MigrationFixDownEvent:
							a.Logger.WriteLine(
								"    %s",
								fmtx.Warning("%s", log.GetMetadata("comment")),
							)
						}
					}
				}
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
			runCmd := run.CommandMigrateUp{
				Driver: cmd.String("driver"),
				Set:    cmd.String("set"),
			}
			out, err := runner.ExecuteMigrateUp(ctx, &runCmd)
			if err != nil {
				return err
			}
			if out.WasPending == 0 {
				a.Logger.WriteSuccess("No pending migrations")
			}

			return nil
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
			runCmd := run.CommandMigrateDown{
				Driver:  cmd.String("driver"),
				Set:     cmd.String("set"),
				Version: cmd.String("version"),
			}
			_, err := runner.ExecuteMigrateDown(ctx, &runCmd)
			if err != nil {
				return err
			}
			a.Logger.WriteSuccess("Migration %s rolled back successfully", runCmd.Version)

			return nil
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
			runCmd := run.CommandMigrateStatus{
				Driver: cmd.String("driver"),
				Set:    cmd.String("set"),
			}
			out, err := runner.ExecuteMigrateStatus(ctx, &runCmd)
			if err != nil {
				return err
			}

			a.Logger.WriteInfo("Migrations status for %s/%s", runCmd.Driver, runCmd.Set)

			for _, row := range out.Rows {
				status := ""
				switch row.Status {
				case "pending":
					status = fmtx.Warning(" %9s ", row.Status)
				case "completed":
					status = fmtx.Success(" %s ", row.Status)
				case "failed":
					status = fmtx.Danger(" %9s ", row.Status)
				}
				a.Logger.WriteLine(
					"%s %-30s %s",
					status,
					row.Version,
					row.Migration,
				)
			}
			return nil
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
			migrateFixGroup,
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

    err := root.Run(context.Background(), args)
	if err != nil {
		a.Logger.WriteError(err.Error())
		return 1
	}

	return 0
}
