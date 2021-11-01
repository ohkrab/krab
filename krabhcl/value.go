package krabhcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Expression struct {
	Expr        hcl.Expression
	EvalContext *hcl.EvalContext
}

func (e Expression) AsBool() bool {
	val, _ := e.Expr.Value(e.EvalContext)
	var boolean bool
	if err := gocty.FromCtyValue(val, &boolean); err == nil {
		return boolean
	}
	return false
}

func (e Expression) Ok() bool {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.IsWhollyKnown() && !val.IsNull()
}

func (e Expression) Type() cty.Type {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.Type()
}
