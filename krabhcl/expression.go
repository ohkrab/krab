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

func (e Expression) TryString() (string, error) {
	val, diags := e.Expr.Value(e.EvalContext)
	if diags.HasErrors() {
		return "", diags.Errs()[0]
	}
	var str string
	err := gocty.FromCtyValue(val, &str)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (e Expression) AsString() string {
	val, _ := e.Expr.Value(e.EvalContext)
	var str string
	if err := gocty.FromCtyValue(val, &str); err == nil {
		return str
	}
	return ""
}

func (e Expression) AsSliceAddr() []*Addr {
	addrs := []*Addr{}
	traversals := e.Expr.Variables()
	for _, t := range traversals {
		addr, _ := ParseTraversalToAddr(t)
		addrs = append(addrs, &addr)
	}
	return addrs
}

func (e Expression) AsSliceString() []string {
	val, _ := e.Expr.Value(e.EvalContext)
	if val.Type().IsTupleType() && val.IsWhollyKnown() {
		vals := val.AsValueSlice()
		ss := []string{}
		for _, v := range vals {
			var str string
			if err := gocty.FromCtyValue(v, &str); err == nil {
				ss = append(ss, str)
			} else {
				return nil
			}
		}
		return ss
	}
	return nil
}

func (e Expression) Ok() bool {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.IsWhollyKnown() && !val.IsNull()
}

func (e Expression) Type() cty.Type {
	val, _ := e.Expr.Value(e.EvalContext)
	return val.Type()
}
