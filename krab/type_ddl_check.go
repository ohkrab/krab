package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLCheck constraint DSL for table DDL.
type DDLCheck struct {
	krabhcl.Source

	Name       string
	Expression string
}

var schemaCheck = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "expression", Required: true},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLCheck) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaCheck)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "expression":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			d.Expression = expr.AsString()

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLCheck) ToSQL(w io.StringWriter) {
	w.WriteString("CONSTRAINT ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" CHECK (")
	w.WriteString(d.Expression)
	w.WriteString(")")
}
