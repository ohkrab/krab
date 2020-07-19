package configs

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type File struct {
	Connections []*Connection

	Variables map[string]cty.Value
	Functions map[string]function.Function
}
