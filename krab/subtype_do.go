package krab

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Do subtype for other types.
type Do struct {
	Action hcl.Expression       `hcl:"action,optional"`
	Inputs map[string]cty.Value `hcl:"inputs,optional"`
	SQL    string               `hcl:"sql,optional"`
}
