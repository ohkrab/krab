package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLCreateTable contains DSL for creating tables.
type DDLCreateTable struct {
	krabhcl.Source

	Name        string
	Unlogged    bool             `hcl:"unlogged,optional"`
	Columns     []*DDLColumn     `hcl:"column,block"`
	PrimaryKeys []*DDLPrimaryKey `hcl:"primary_key,block"`
	ForeignKeys []*DDLForeignKey `hcl:"foreign_key,block"`
	Uniques     []*DDLUnique     `hcl:"unique,block"`
	Checks      []*DDLCheck      `hcl:"check,block"`
}

var schemaCreateTable = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "column",
			LabelNames: []string{"name", "type"},
		},
		{
			Type:       "primary_key",
			LabelNames: []string{},
		},
		{
			Type:       "foreign_key",
			LabelNames: []string{},
		},
		{
			Type:       "unique",
			LabelNames: []string{},
		},
		{
			Type:       "check",
			LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{Name: "unlogged", Required: false},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLCreateTable) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Name = block.Labels[0]
	d.Columns = []*DDLColumn{}
	d.PrimaryKeys = []*DDLPrimaryKey{}
	d.ForeignKeys = []*DDLForeignKey{}
	d.Uniques = []*DDLUnique{}
	d.Checks = []*DDLCheck{}

	content, diags := block.Body.Content(schemaCreateTable)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "column":
			column := new(DDLColumn)
			err := column.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			d.Columns = append(d.Columns, column)

		case "primary_key":
			pk := new(DDLPrimaryKey)
			err := pk.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			d.PrimaryKeys = append(d.PrimaryKeys, pk)

		case "foreign_key":
			fk := new(DDLForeignKey)
			err := fk.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			d.ForeignKeys = append(d.ForeignKeys, fk)

		case "unique":
			unique := new(DDLUnique)
			err := unique.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			d.Uniques = append(d.Uniques, unique)

		case "check":
			check := new(DDLCheck)
			err := check.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			d.Checks = append(d.Checks, check)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "unlogged":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			d.Unlogged = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLCreateTable) ToSQL(w io.StringWriter) {
	w.WriteString("CREATE")
	if d.Unlogged {
		w.WriteString(" UNLOGGED")
	}
	w.WriteString(" TABLE ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString("(\n")

	hasPK := len(d.PrimaryKeys) > 0
	hasFK := len(d.ForeignKeys) > 0
	hasUnique := len(d.Uniques) > 0
	hasCheck := len(d.Checks) > 0

	for i, col := range d.Columns {
		w.WriteString("  ")
		col.ToSQL(w)
		if i < len(d.Columns)-1 {
			w.WriteString(",")
			w.WriteString("\n")
		}
	}

	if hasPK {
		for _, pk := range d.PrimaryKeys {
			w.WriteString("\n, ")
			pk.ToSQL(w)
		}
	}
	if hasFK {
		for _, fk := range d.ForeignKeys {
			w.WriteString("\n, ")
			fk.ToSQL(w)
		}
	}
	if hasUnique {
		for _, u := range d.Uniques {
			w.WriteString("\n, ")
			u.ToSQL(w)
		}
	}
	if hasCheck {
		for _, c := range d.Checks {
			w.WriteString("\n, ")
			c.ToSQL(w)
		}
	}

	w.WriteString("\n)")
}
