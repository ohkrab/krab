package krab

import (
	"github.com/ohkrab/krab/configs"
	"github.com/zclconf/go-cty/cty"
)

type Evaluator struct {
	Config         *configs.Config
	VariableValues map[string]map[string]cty.Value
}
