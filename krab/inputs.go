package krab

import "github.com/zclconf/go-cty/cty"

// Inputs are params passed to command.
type Inputs map[string]interface{}

func InputsFromCtyInputs(vals map[string]cty.Value) Inputs {
	inputs := Inputs{}
	for k, v := range vals {
		str := v.AsString()
		inputs[k] = str
	}

	return inputs
}

func (i Inputs) Merge(other Inputs) {
	for k, v := range other {
		i[k] = v
	}
}
