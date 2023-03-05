package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// TestExampleIt represents one use case for test example that contain queries and assertions.
type TestExampleIt struct {
	krabhcl.Source

	Name    string
	Queries []*TestQuery
}

var schemaTestExampleIt = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "query",
			LabelNames: []string{"sql"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (it *TestExampleIt) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	it.Source.Extract(block)

	it.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaTestExampleIt)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "query":
			q := new(TestQuery)
			err := q.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			it.Queries = append(it.Queries, q)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}
