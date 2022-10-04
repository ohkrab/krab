package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLCreateIndex contains DSL for creating indicies.
type DDLCreateIndex struct {
	krabhcl.Source

	Table        string
	Name         string
	Unique       bool
	Concurrently bool
	Columns      []string
	Include      []string
	Using        string
	Where        string
}

var schemaCreateIndex = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "column",
			LabelNames: []string{"name", "type"},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
		{Name: "include", Required: false},
		{Name: "unique", Required: false},
		{Name: "using", Required: false},
		{Name: "where", Required: false},
		{Name: "concurrently", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLCreateIndex) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Table = block.Labels[0]
	d.Name = block.Labels[1]
	d.Unique = false
	d.Concurrently = false
	d.Columns = []string{}
	d.Include = []string{}
	d.Using = ""
	d.Where = ""

	content, diags := block.Body.Content(schemaCreateIndex)
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

		case "unique":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Unique = val

		case "using":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			d.Using = val

		case "where":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			d.Where = val

		case "concurrently":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Concurrently = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLCreateIndex) ToSQL(w io.StringWriter) {
	w.WriteString("CREATE")
	if d.Unique {
		w.WriteString(" UNIQUE")
	}
	w.WriteString(" INDEX")
	if d.Concurrently {
		w.WriteString(" CONCURRENTLY")
	}
	w.WriteString(" ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" ON ")
	w.WriteString(krabdb.QuoteIdent(d.Table))
	if d.Using != "" {
		w.WriteString(" USING ")
		w.WriteString(d.Using)
	}
	w.WriteString(" (")
	for i, col := range d.Columns {
		w.WriteString(krabdb.QuoteIdent(col))
		if i < len(d.Columns)-1 {
			w.WriteString(",")
		}
	}
	w.WriteString(")")

	if len(d.Include) > 0 {
		w.WriteString(" INCLUDE (")
		for i, col := range d.Include {
			w.WriteString(krabdb.QuoteIdent(col))
			if i < len(d.Columns)-1 {
				w.WriteString(",")
			}
		}
		w.WriteString(")")
	}

	if d.Where != "" {
		w.WriteString(" WHERE (")
		w.WriteString(d.Where)
		w.WriteString(")")
	}
}
