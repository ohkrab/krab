package krabhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Expression struct {
	Expr        hcl.Expression
	EvalContext *hcl.EvalContext
}

func (e Expression) Addr() (Addr, error) {
	traversals := e.Expr.Variables()
	if len(traversals) != 1 {
		return Addr{}, fmt.Errorf("Failed to extract single addr from HCL expression")
	}

	t := traversals[0]
	parsedAddr, err := ParseTraversalToAddr(t)
	if err != nil {
		return Addr{}, err
	}
	return parsedAddr, nil
}

func (e Expression) AsBool() bool {
	val, _ := e.Expr.Value(e.EvalContext)
	var boolean bool
	if err := gocty.FromCtyValue(val, &boolean); err == nil {
		return boolean
	}
	return false
}

func (e Expression) AsInt64() int64 {
	val, _ := e.Expr.Value(e.EvalContext)
	var number int64
	if err := gocty.FromCtyValue(val, &number); err == nil {
		return number
	}
	return 0
}

func (e Expression) AsFloat64() float64 {
	val, _ := e.Expr.Value(e.EvalContext)
	var number float64
	if err := gocty.FromCtyValue(val, &number); err == nil {
		return number
	}
	return 0
}

func (e Expression) AsString() string {
	val, _ := e.Expr.Value(e.EvalContext)
	var str string
	if err := gocty.FromCtyValue(val, &str); err == nil {
		return str
	}
	return ""
}

func (e Expression) Ok() bool {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.IsWhollyKnown() && !val.IsNull()
}

func (e Expression) Type() cty.Type {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.Type()
}
