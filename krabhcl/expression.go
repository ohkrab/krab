package krabhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
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

func (e Expression) Bool() (bool, error) {
	val, diags := e.Expr.Value(e.EvalContext)
	if diags.HasErrors() {
		return false, diags.Errs()[0]
	}
	var boolean bool
	if err := gocty.FromCtyValue(val, &boolean); err != nil {
		return false, err
	}
	return boolean, nil
}

func (e Expression) Int64() (int64, error) {
	val, diags := e.Expr.Value(e.EvalContext)
	if diags.HasErrors() {
		return 0, diags.Errs()[0]
	}
	var number int64
	if err := gocty.FromCtyValue(val, &number); err != nil {
		return 0, err
	}
	return number, nil
}

func (e Expression) AsFloat64() (float64, error) {
	val, diags := e.Expr.Value(e.EvalContext)
	if diags.HasErrors() {
		return 0, diags.Errs()[0]
	}
	var number float64
	if err := gocty.FromCtyValue(val, &number); err != nil {
		return 0, err
	}
	return number, nil
}

func (e Expression) String() (string, error) {
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

func (e Expression) SliceAddr() ([]*Addr, error) {
	addrs := []*Addr{}
	traversals := e.Expr.Variables()
	for _, t := range traversals {
		addr, err := ParseTraversalToAddr(t)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, &addr)
	}
	return addrs, nil
}

func (e Expression) SliceString() ([]string, error) {
	val, diags := e.Expr.Value(e.EvalContext)
	if diags.HasErrors() {
		return nil, diags.Errs()[0]
	}
	if val.Type().IsTupleType() && val.IsWhollyKnown() {
		vals := val.AsValueSlice()
		ss := []string{}
		for _, v := range vals {
			var str string
			if err := gocty.FromCtyValue(v, &str); err == nil {
				ss = append(ss, str)
			} else {
				return nil, err
			}
		}
		return ss, nil
	}
	return nil, fmt.Errorf("Inwalid types in a list, expected strings")
}
