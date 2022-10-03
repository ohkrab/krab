package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLGeneratedColumn DSL.
type DDLGeneratedColumn struct {
	As string
}

var schemaGeneratedColumn = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "as", Required: true},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLGeneratedColumn) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	content, diags := block.Body.Content(schemaGeneratedColumn)
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
		case "as":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			d.As = expr.AsString()

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLGeneratedColumn) ToSQL(w io.StringWriter) {
	w.WriteString("GENERATED ALWAYS AS (")
	w.WriteString(d.As)
	w.WriteString(") STORED")
}
