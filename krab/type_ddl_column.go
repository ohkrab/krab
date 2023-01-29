package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLColumn DSL for table DDL.
type DDLColumn struct {
	krabhcl.Source

	Name      string
	Type      string
	Null      bool
	Default   string
	Identity  *DDLIdentity
	Generated *DDLGeneratedColumn
}

var schemaColumn = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "generated",
			LabelNames: []string{},
		},
		{
			Type:       "identity",
			LabelNames: []string{},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{Name: "null", Required: false},
		{Name: "default", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLColumn) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Name = block.Labels[0]
	d.Type = block.Labels[1]
	d.Null = true
	d.Identity = nil
	d.Generated = nil

	content, diags := block.Body.Content(schemaColumn)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "identity":
			d.Identity = &DDLIdentity{}
			err := d.Identity.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}

		case "generated":
			d.Generated = &DDLGeneratedColumn{}
			err := d.Generated.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "null":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Null = val

		case "default":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			d.Default = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLColumn) ToSQL(w io.StringWriter) {
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" ")
	w.WriteString(d.Type)

	if !d.Null {
		w.WriteString(" NOT NULL")
	}

	if d.Identity != nil {
		w.WriteString(" ")
		d.Identity.ToSQL(w)
	}

	if d.Generated != nil {
		w.WriteString(" ")
		d.Generated.ToSQL(w)
	}

	if d.Default != "" {
		w.WriteString(" DEFAULT ")
		w.WriteString(d.Default)
	}
}

// ToKCL converts migration definition to KCL.
func (d *DDLColumn) ToKCL(w io.StringWriter) {
	w.WriteString("column ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" ")
	w.WriteString(krabdb.QuoteIdent(d.Type))
	w.WriteString(" {")
	if d.Identity != nil {
		w.WriteString("\n  identity {}\n")
	}
	if !d.Null {
		w.WriteString("\n  null = false\n")
	}
	w.WriteString("}")
}
