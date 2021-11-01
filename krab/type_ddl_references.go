package krab

import (
	"io"
	"strings"

	"github.com/ohkrab/krab/krabdb"
)

// DDLReferences DSL for ForeignKey.
type DDLReferences struct {
	Table    string   `hcl:"table,label"`
	Columns  []string `hcl:"columns"`
	OnDelete string   `hcl:"on_delete,optional"`
	OnUpdate string   `hcl:"on_update,optional"`
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
