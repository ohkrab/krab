package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLIdentity DSL.
type DDLIdentity struct {
	krabhcl.Source
}

var schemaIdentity = &hcl.BodySchema{
	Blocks:     []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{},
}

// DecodeHCL parses HCL into struct.
func (d *DDLIdentity) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	content, diags := block.Body.Content(schemaColumn)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, _ := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLIdentity) ToSQL(w io.StringWriter) {
	w.WriteString("GENERATED ALWAYS AS IDENTITY")
}
