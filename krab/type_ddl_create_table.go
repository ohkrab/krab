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

	Name        string           `hcl:"name,label"`
	Unlogged    bool             `hcl:"unlogged,optional"`
	Columns     []*DDLColumn     `hcl:"column,block"`
	PrimaryKeys []*DDLPrimaryKey `hcl:"primary_key,block"`
	ForeignKeys []*DDLForeignKey `hcl:"foreign_key,block"`
	Uniques     []*DDLUnique     `hcl:"unique,block"`
	Checks      []*DDLCheck      `hcl:"check,block"`
}

var DDLCreateTableSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "create_table",
			LabelNames: []string{"name"},
		},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLCreateTable) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	d.Source.Extract(block)

	d.Columns = []*DDLColumn{}
	d.PrimaryKeys = []*DDLPrimaryKey{}
	d.ForeignKeys = []*DDLForeignKey{}
	d.Uniques = []*DDLUnique{}
	d.Checks = []*DDLCheck{}

	content, diags := block.Body.Content(DDLCreateTableSchema)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `table` block: %s", diags.Error())
	}

	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `table` attributes: %s", diags.Error())
	}

	for k, v := range attrs {
		fmt.Println(k, v)
	}

	for _, b := range content.Blocks {
		switch b.Type {
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
