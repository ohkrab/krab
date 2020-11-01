package krab

import (
	"github.com/ohkrab/krab/configs"
	"github.com/ohkrab/krab/diagnostics"
	"github.com/zclconf/go-cty/cty/function"
)

type Context struct {
	config *configs.Config
}

type ContextOpts struct {
	Config    *configs.Config
	Functions map[string]function.Function
}

func NewContext(opts *ContextOpts) (*Context, diagnostics.List) {
	return &Context{
		config: opts.Config,
	}, nil
}

func (c *Context) Graph() (*Graph, diagnostics.List) {
	steps := []GraphTransformer{
		&ReferenceTransformer{},
	}
	builder := &GraphBuilder{Steps: steps}
	return builder.Build(c.config.Module)
}

func (c *Context) Eval() diagnostics.List {
	graph, diags := c.Graph()
	walker := &ContextGraphWalker{
		Context: c,
	}
	graph.Walk(walker)
	return diags
}
