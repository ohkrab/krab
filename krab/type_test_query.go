package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

type TestQuery struct {
	krabhcl.Source

	Query string
	Rows  []*TestQueryRow
}

var schemaTestQuery = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "row",
			LabelNames: []string{"scope"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (q *TestQuery) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	q.Source.Extract(block)
	q.Query = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQuery)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "row":
			row := new(TestQueryRow)
			err := row.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			q.Rows = append(q.Rows, row)

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
