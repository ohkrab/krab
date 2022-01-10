package krab

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Do subtype for other types.
type Do struct {
	MigrationSet hcl.Expression       `hcl:"migration_set,optional"`
	CtyInputs    map[string]cty.Value `hcl:"inputs,optional"`
	SQL          string               `hcl:"sql,optional"`
}

func (d *Do) Inputs() Inputs {
	inputs := Inputs{}
	for k, v := range d.CtyInputs {
		str := v.AsString()
		inputs[k] = str
	}

	return inputs
}
