package krab

import (
	"github.com/ohkrab/krab/configs"
	"github.com/ohkrab/krab/diagnostics"
)

type Context struct {
	config *configs.Config
}

type ContextOpts struct {
	Config *configs.Config
}

func NewContext(opts *ContextOpts) (*Context, diagnostics.List) {
	return &Context{
		config: opts.Config,
	}, nil
}

func (ctx *Context) Graph() (*Graph, diagnostics.List) {
	steps := []GraphTransformer{}
	builder := &GraphBuilder{Steps: steps}
	return builder.Build(ctx.config.Module)
}
