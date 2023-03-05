package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

type TestQueryCol struct {
	krabhcl.Source

	Message string
	Assert  string
}

var schemaTestQueryCol = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "assert",
			Required: true,
		},
	},
}

func (col *TestQueryCol) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	col.Source.Extract(block)
	col.Message = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQueryCol)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {

		case "assert":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			col.Assert = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}
