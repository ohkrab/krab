package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLDropTable contains DSL for dropping tables.
type DDLDropTable struct {
	krabhcl.Source

	Name string
}

var schemaDropTable = &hcl.BodySchema{}

// DecodeHCL parses HCL into struct.
func (d *DDLDropTable) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaDropTable)
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
func (d *DDLDropTable) ToSQL(w io.StringWriter) {
	w.WriteString("DROP TABLE ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
}
