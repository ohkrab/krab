package krab

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLForeignKey constraint DSL for table DDL.
type DDLForeignKey struct {
	krabhcl.Source

	Columns    []string
	References DDLReferences
}

var schemaForeignKey = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "references", LabelNames: []string{"table"}},
	},
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLForeignKey) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Columns = []string{}

	content, diags := block.Body.Content(schemaForeignKey)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "references":
			err := d.References.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "columns":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			columns := expr.AsSliceString()
			d.Columns = append(d.Columns, columns...)

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLForeignKey) ToSQL(w io.StringWriter) {
	w.WriteString("FOREIGN KEY (")
	cols := krabdb.QuoteIdentStrings(d.Columns)
	w.WriteString(strings.Join(cols, ","))
	w.WriteString(") ")
	d.References.ToSQL(w)
}
