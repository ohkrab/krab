package krab

import (
	"context"
)

// Cmd is a command that app can execute.
type Cmd interface {
	// Name that is mounted at API path or CLI.
	Name() []string

	// HttpMethod that is used for API call.
	HttpMethod() string

	// Do executes the action.
	Do(ctx context.Context, opts CmdOpts) (interface{}, error)
}

// CmdOpts are options passed to command.
type CmdOpts struct {
	Inputs
}
