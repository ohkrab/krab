package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// TestExample represents test runner configuration.
type TestExample struct {
	krabhcl.Source

	// Set              *SetRuntimeParameters `hcl:"set,block"`
	Name string
	Its  []*TestExampleIt
	Xits []*TestExampleIt
}

var schemaTestExample = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "it",
			LabelNames: []string{"name"},
		},
		{
			Type:       "xit",
			LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (t *TestExample) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "test", Labels: []string{t.Name}}
}

func (t *TestExample) Validate() error {
	return ErrorCoalesce()
}

func (t *TestExample) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	t.Source.Extract(block)
	t.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaTestExample)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "it":
			it := new(TestExampleIt)
			err := it.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			t.Its = append(t.Its, it)

		case "xit":
			it := new(TestExampleIt)
			err := it.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			t.Xits = append(t.Xits, it)

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
