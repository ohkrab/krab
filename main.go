package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krab"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		PartsOrder: []string{zerolog.LevelFieldName, zerolog.CallerFieldName, zerolog.MessageFieldName},
	}
	log := zerolog.New(output).With().Logger()

	dir, err := optGetDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read config dir")
		os.Exit(1)
	}

	parser := krab.NewParser()
	config, err := parser.LoadConfigDir(dir)
	if err != nil {
		log.Error().Err(err).Msg("Parser error")
		os.Exit(1)
	}

	cmds := make([]*cli.Command, 1, 10)
	cmds[0] = &cli.Command{
		Name:  "version",
		Usage: "Print version",
		Action: func(c *cli.Context) error {
			fmt.Fprintln(c.App.Writer, krab.InfoName, krab.InfoVersion)
			fmt.Fprintln(c.App.Writer, "Build", krab.InfoCommit, krab.InfoBuildDate)
			return nil
		},
	}

	migrateUpCmds := make([]*cli.Command, 0, len(config.MigrationSets))
	migrateDownCmds := make([]*cli.Command, 0, len(config.MigrationSets))
	for _, set := range config.MigrationSets {
		migrateUpCmds = append(migrateUpCmds, migrateUpCommand(ctx, set.RefName, config))
		migrateDownCmds = append(migrateDownCmds, migrateDownCommand(ctx, set.RefName, config))
	}
	cmdUp := &cli.Command{Name: "up", Subcommands: migrateUpCmds}
	cmdDown := &cli.Command{Name: "down", Subcommands: migrateDownCmds}
	cmds = append(cmds, &cli.Command{
		Name:        "migrate",
		Usage:       "Migration commands",
		Subcommands: []*cli.Command{cmdUp, cmdDown},
	})

	app := &cli.App{
		Name:     "krab",
		Usage:    fmt.Sprint("PostgreSQL tool 🐘\n   ", krab.InfoWWW),
		Commands: cmds,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("Action failed")
	}
}

func withPg(f func(db *sqlx.DB) error) error {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	return f(db)
}

func optGetDir() (string, error) {
	if dir := os.Getenv("KRAB_DIR"); dir != "" {
		return dir, nil
	}

	return os.Getwd()
}

func migrateUpCommand(ctx context.Context, name string, config *krab.Config) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: fmt.Sprintf("Migrate `%s` up", name),
		Action: func(c *cli.Context) error {
			action := krab.ActionMigrateUp{Set: config.MigrationSets[name]}
			return withPg(func(db *sqlx.DB) error {
				return action.Run(ctx, db)
			})
		},
	}
}

func migrateDownCommand(ctx context.Context, name string, config *krab.Config) *cli.Command {
	return &cli.Command{
		Name:            name,
		Usage:           fmt.Sprintf("Migrate `%s` down", name),
		ArgsUsage:       "[version]",
		HideHelp:        true,
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("Requires [version] argument to be specified")
			}

			action := krab.ActionMigrateDown{
				Set:           config.MigrationSets[name],
				DownMigration: krab.SchemaMigration{c.Args().First()},
			}
			return withPg(func(db *sqlx.DB) error {
				return action.Run(ctx, db)
			})
		},
	}
}
