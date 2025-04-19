package run

import (
	"context"

	"github.com/ohkrab/krab/ferro/config"
)

type Initializer struct {
	fs *config.Filesystem
}

func NewInitializer(fs *config.Filesystem) *Initializer {
	return &Initializer{fs: fs}
}

type InitializeOptions struct {
}

func (init *Initializer) Initialize(ctx context.Context, opts InitializeOptions) error {
	return nil
}
