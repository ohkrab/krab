package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

type TestQueryRow struct {
	krabhcl.Source

	Scope string
	Cols  []*TestQueryCol
}

var schemaTestQueryRow = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "col",
			LabelNames: []string{"message"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (row *TestQueryRow) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	row.Source.Extract(block)
	row.Scope = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQueryRow)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "col":
			col := new(TestQueryCol)
			err := col.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			row.Cols = append(row.Cols, col)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, _ := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}
