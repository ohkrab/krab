package krab

import (
	"context"

	"github.com/ohkrab/krab/krabhcl"
)

// Cmd is a command that app can execute.
type Cmd interface {
	// Addr associated with name
	Addr() krabhcl.Addr

	// Name that is mounted at API path or CLI.
	Name() []string

	// HttpMethod that is used for API call.
	HttpMethod() string

	// Do executes the action.
	Do(ctx context.Context, opts CmdOpts) (any, error)
}

// CmdOpts are options passed to command.
type CmdOpts struct {
	NamedInputs
	PositionalInputs
}
