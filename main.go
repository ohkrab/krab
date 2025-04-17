package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	_ "embed"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/parser"
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

type hexWriter struct{}

func (w *hexWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		fmt.Printf("%x", b)
	}
	fmt.Printf("\n")

	return len(p), nil
}

type genericType struct {
	s string
}

func (g *genericType) Set(value string) error {
	g.s = value
	return nil
}

func (g *genericType) String() string {
	return g.s
}

func main() {
	filesystem := config.NewFilesystem(".")
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

	for _, migrationSet := range cfg.MigrationSets {
		fmt.Printf("MigrationSet: %s\n", migrationSet.Metadata.Name)
	}
	root := &cli.Command{
		Name:      "ferro",
		Version:   "v0.20.0",
		Copyright: "(c) @qbart",
		Usage:     "ferroDB is a tool for managing your databases",
		UsageText: "ferro [global options] command [command options] [arguments...]",
		ArgsUsage: "[args and such]",
	}
	root.Run(context.Background(), os.Args)
	return
	cmd := &cli.Command{
		Name:    "kənˈtrīv",
		Version: "v19.99.0",
		/*Authors: []any{
		    &cli.Author{
		        Name:  "Example Human",
		        Email: "human@example.com",
		    },
		},*/
		Copyright: "(c) 1999 Serious Enterprise",
		Usage:     "demonstrate available API",
		UsageText: "contrive - demonstrating the available API",
		ArgsUsage: "[args and such]",
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "doo",
				Aliases:     []string{"do"},
				Category:    "motion",
				Usage:       "do the doo",
				UsageText:   "doo - does the dooing",
				Description: "no really, there is a lot of dooing to be done",
				ArgsUsage:   "[arrgh]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				},
				Commands: []*cli.Command{
					&cli.Command{
						Name:   "wop",
						Action: wopAction,
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				ShellComplete: func(ctx context.Context, cmd *cli.Command) {
					fmt.Fprintf(cmd.Root().Writer, "--better\n")
				},
				Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
					fmt.Fprintf(cmd.Root().Writer, "brace for impact\n")
					return nil, nil
				},
				After: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Fprintf(cmd.Root().Writer, "did we lose anyone?\n")
					return nil
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cmd.FullName()
					cmd.HasName("wop")
					cmd.Names()
					cmd.VisibleFlags()
					fmt.Fprintf(cmd.Root().Writer, "dodododododoodododddooooododododooo\n")
					if cmd.Bool("forever") {
						cmd.Run(ctx, nil)
					}
					return nil
				},
				OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
					fmt.Fprintf(cmd.Root().Writer, "for shame\n")
					return err
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "fancy"},
			&cli.BoolFlag{Value: true, Name: "fancier"},
			&cli.DurationFlag{Name: "howlong", Aliases: []string{"H"}, Value: time.Second * 3},
			&cli.FloatFlag{Name: "howmuch"},
			&cli.IntFlag{Name: "longdistance", Validator: func(t int64) error {
				if t < 10 {
					return fmt.Errorf("10 miles isnt long distance!!!!")
				}
				return nil
			}},
			&cli.IntSliceFlag{Name: "intervals"},
			&cli.StringFlag{Name: "dance-move", Aliases: []string{"d"}, Validator: func(move string) error {
				moves := []string{"salsa", "tap", "two-step", "lock-step"}
				if !slices.Contains(moves, move) {
					return fmt.Errorf("Havent learnt %s move yet", move)
				}
				return nil
			}},
			&cli.StringSliceFlag{Name: "names", Aliases: []string{"N"}},
			&cli.UintFlag{Name: "age"},
		},
		EnableShellCompletion: true,
		HideHelp:              false,
		HideVersion:           false,
		ShellComplete: func(ctx context.Context, cmd *cli.Command) {
			fmt.Fprintf(cmd.Root().Writer, "lipstick\nkiss\nme\nlipstick\nringo\n")
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			fmt.Fprintf(cmd.Root().Writer, "HEEEERE GOES\n")
			return nil, nil
		},
		After: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Root().Writer, "Phew!\n")
			return nil
		},
		CommandNotFound: func(ctx context.Context, cmd *cli.Command, command string) {
			fmt.Fprintf(cmd.Root().Writer, "Thar be no %q here.\n", command)
		},
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cmd.Root().Writer, "WRONG: %#v\n", err)
			return nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cli.DefaultAppComplete(ctx, cmd)
			cli.HandleExitCoder(errors.New("not an exit coder, though"))
			cli.ShowAppHelp(cmd)
			cli.ShowCommandHelp(ctx, cmd, "also-nope")
			cli.ShowSubcommandHelp(cmd)
			cli.ShowVersion(cmd)

			fmt.Printf("%#v\n", cmd.Root().Command("doo"))
			if cmd.Bool("infinite") {
				cmd.Root().Run(ctx, []string{"app", "doo", "wop"})
			}

			if cmd.Bool("forevar") {
				cmd.Root().Run(ctx, nil)
			}
			fmt.Printf("%#v\n", cmd.Root().VisibleCategories())
			fmt.Printf("%#v\n", cmd.Root().VisibleCommands())
			fmt.Printf("%#v\n", cmd.Root().VisibleFlags())

			fmt.Printf("%#v\n", cmd.Args().First())
			if cmd.Args().Len() > 0 {
				fmt.Printf("%#v\n", cmd.Args().Get(1))
			}
			fmt.Printf("%#v\n", cmd.Args().Present())
			fmt.Printf("%#v\n", cmd.Args().Tail())

			ec := cli.Exit("ohwell", 86)
			fmt.Fprintf(cmd.Root().Writer, "%d", ec.ExitCode())
			fmt.Printf("made it!\n")
			return ec
		},
		Metadata: map[string]interface{}{
			"layers":          "many",
			"explicable":      false,
			"whatever-values": 19.99,
		},
	}

	if os.Getenv("HEXY") != "" {
		cmd.Writer = &hexWriter{}
		cmd.ErrWriter = &hexWriter{}
	}

	cmd.Run(context.Background(), os.Args)
}

func wopAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Fprintf(cmd.Root().Writer, ":wave: over here, eh\n")
	return nil
}
