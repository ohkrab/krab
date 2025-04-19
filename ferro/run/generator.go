package run

import (
	"context"

	"github.com/ohkrab/krab/ferro/config"
)

type Generator struct {
	fs *config.Filesystem
}

func NewGenerator(fs *config.Filesystem) *Generator {
	return &Generator{fs: fs}
}

type GenerateOptions struct {
}

func (g *Generator) Generate(ctx context.Context, opts GenerateOptions) error {
	return nil
}
