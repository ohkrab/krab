package krab

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLPrimaryKey constraint DSL for table DDL.
type DDLPrimaryKey struct {
	krabhcl.Source

	Columns []string
	Include []string
}

var schemaPrimaryKey = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
		{Name: "include", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLPrimaryKey) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Columns = []string{}
	d.Include = []string{}

	content, diags := block.Body.Content(schemaPrimaryKey)
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

		case "include":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.SliceString()
			if err != nil {
				return err
			}
			d.Include = append(d.Include, val...)

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLPrimaryKey) ToSQL(w io.StringWriter) {
	w.WriteString("PRIMARY KEY (")
	cols := krabdb.QuoteIdentStrings(d.Columns)
	w.WriteString(strings.Join(cols, ","))
	w.WriteString(")")

	if len(d.Include) > 0 {
		w.WriteString(" INCLUDE (")
		include := krabdb.QuoteIdentStrings(d.Include)
		w.WriteString(strings.Join(include, ","))
		w.WriteString(")")
	}
}

// ToKCL converts migration definition to KCL.
func (d *DDLPrimaryKey) ToKCL(w io.StringWriter) {
	w.WriteString("      primary_key {\n")
	w.WriteString("        columns = [")
	w.WriteString(strings.Join(krabdb.QuoteIdentStrings(d.Columns), ""))
	w.WriteString("]\n")
	w.WriteString("      }\n")
}
