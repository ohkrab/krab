package krab

import (
	"github.com/ohkrab/krab/configs"
	"github.com/ohkrab/krab/diagnostics"
)

func Load(dir string) (*Context, diagnostics.List) {
	diags := diagnostics.New()
	parser := configs.NewParser()
	config, hclDiags := parser.LoadConfigDir(dir)
	diags.Append(hclDiags)

	if diags.HasErrors() {
		return nil, diags
	}

	ctx, ctxDiags := NewContext(&ContextOpts{
		Config: config,
	})
	diags.Append(ctxDiags)

	return ctx, diags
}
