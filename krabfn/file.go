package krabfn

import (
	"github.com/spf13/afero"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// FnFileRead reads the whole file or returns error.
// https://github.com/hashicorp/hcl/blob/main/guide/go_expression_eval.rst
var FnFileRead = func(fs afero.Afero) function.Function {
	return function.New(
		&function.Spec{
			Params: []function.Parameter{
				{Name: "path", Type: cty.String},
			},
			Type: function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				path := args[0].AsString()
				content, err := fs.ReadFile(path)
				if err != nil {
					return cty.NilVal, err
				}
				return cty.StringVal(string(content)), nil
			},
		},
	)
}
