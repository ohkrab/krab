package krabfn

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/spf13/afero"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// EvalContext builds default EvalContext for hcl parser.
func EvalContext(fs afero.Afero) *hcl.EvalContext {
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: map[string]function.Function{
			"file_read": FnFileRead(fs),
		},
	}

	return ctx
}
