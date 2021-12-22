package krabcli

import (
	"fmt"

	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabapi"
)

// CmdAgent starts a web server
type CmdAgent struct {
	Registry *krab.CmdRegistry
}

func (a *CmdAgent) Help() string {
	return fmt.Sprint(
		`Usage: krab agent`,
		"\n\n",
		` 
Starts an HTTP server.
`,
	)
}

func (a *CmdAgent) Synopsis() string {
	return fmt.Sprintf("HTTP API")
}

// Run in CLI.
func (a *CmdAgent) Run(args []string) int {
	agent := krabapi.Agent{Registry: a.Registry}
	err := agent.Run()
	if err != nil {
		return 1
	}
	return 0
}
