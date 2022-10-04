package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLDropIndex contains DSL for dropping indicies.
type DDLDropIndex struct {
	krabhcl.Source

	Name         string
	Cascade      bool
	Concurrently bool
}

var schemaDropIndex = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "cascade", Required: false},
		{Name: "concurrently", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLDropIndex) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Name = block.Labels[0]
	d.Cascade = false
	d.Concurrently = false

	content, diags := block.Body.Content(schemaDropIndex)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for k, v := range content.Attributes {
		switch k {
		case "concurrently":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Concurrently = val

		case "cascade":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Cascade = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLDropIndex) ToSQL(w io.StringWriter) {
	w.WriteString("DROP INDEX")
	if d.Concurrently {
		w.WriteString(" CONCURRENTLY")
	}
	w.WriteString(" ")
	w.WriteString(krabdb.QuoteIdentWithDots(d.Name))
	if d.Cascade {
		w.WriteString(" CASCADE")
	}
}
