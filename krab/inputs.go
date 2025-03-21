package krab

import "github.com/zclconf/go-cty/cty"

// NamedInputs are params passed to command.
type NamedInputs map[string]any

// Inputs are params passed to command.
type PositionalInputs []string

func InputsFromCtyInputs(vals map[string]cty.Value) NamedInputs {
	inputs := NamedInputs{}
	for k, v := range vals {
		str := v.AsString()
		inputs[k] = str
	}

	return inputs
}

func (i NamedInputs) Merge(other NamedInputs) {
	for k, v := range other {
		i[k] = v
	}
}
