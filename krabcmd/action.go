package krabcmd

import (
	"context"
	"io"
)

type Action interface {
	Command() []string
	HttpEndpoint() (method string, path string)
	Do(ctx context.Context, w io.Writer) error
}
