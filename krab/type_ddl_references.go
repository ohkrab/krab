package krab

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLReferences DSL for ForeignKey.
type DDLReferences struct {
	krabhcl.Source

	Table    string
	Columns  []string
	OnDelete string
	OnUpdate string
}

var schemaForeignKeyReferences = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
		{Name: "on_delete", Required: false},
		{Name: "on_update", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLReferences) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Columns = []string{}
	d.Table = block.Labels[0]

	content, diags := block.Body.Content(schemaForeignKeyReferences)
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
		case "columns":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.SliceString()
			if err != nil {
				return err
			}
			d.Columns = append(d.Columns, val...)

		case "on_delete":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			d.OnDelete = val

		case "on_update":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			d.OnUpdate = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLReferences) ToSQL(w io.StringWriter) {
	w.WriteString("REFERENCES ")
	w.WriteString(krabdb.QuoteIdent(d.Table))
	w.WriteString("(")
	cols := krabdb.QuoteIdentStrings(d.Columns)
	w.WriteString(strings.Join(cols, ","))
	w.WriteString(")")

	if d.OnDelete != "" {
		w.WriteString(" ON DELETE ")
		w.WriteString(d.OnDelete)
	}

	if d.OnUpdate != "" {
		w.WriteString(" ON UPDATE ")
		w.WriteString(d.OnUpdate)
	}
}
