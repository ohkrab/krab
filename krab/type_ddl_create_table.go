package krab

import (
	"io"

	"github.com/ohkrab/krab/krabdb"
)

// DDLCreateTable contains DSL for creating tables.
type DDLCreateTable struct {
	Name        string           `hcl:"name,label"`
	Unlogged    bool             `hcl:"unlogged,optional"`
	Columns     []*DDLColumn     `hcl:"column,block"`
	PrimaryKeys []*DDLPrimaryKey `hcl:"primary_key,block"`
	ForeignKeys []*DDLForeignKey `hcl:"foreign_key,block"`
	Uniques     []*DDLUnique     `hcl:"unique,block"`
	Checks      []*DDLCheck      `hcl:"check,block"`
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
