package krab

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

// App data
type App struct {
}

// New creates krab instance
func New(dir string) *App {
	return &App{}
}

// Init downloads plugins to cache dir.
func (a *App) Init() {
	fmt.Println("TODO: init")
}

// Agent starts agent mode.
func (a *App) Agent() {
	agent := Agent{}
	agent.Run()
}

// Run starts the app and registers all plugins
func (a *App) Run() {
	// root
	rootCmd := &cobra.Command{
		Use:   "krab",
		Short: "Krab is a pluggable database/automation tool",
	}

	// version cmd
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Krab version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Krab v", Version, "\n")
		},
	}

	// agent cmd
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Start agent",
		Run: func(cmd *cobra.Command, args []string) {
			a.Agent()
		},
	}

	// register internal commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(agentCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
