package krab

import (
	"io"

	"github.com/ohkrab/krab/krabdb"
)

// DDLCreateIndex contains DSL for creating indicies.
type DDLCreateIndex struct {
	Table string `hcl:"table,label"`
	Name  string `hcl:"name,label"`

	Unique       bool     `hcl:"unique,optional"`
	Concurrently bool     `hcl:"concurrently,optional"`
	Columns      []string `hcl:"columns"`
	Include      []string `hcl:"include,optional"`
	Using        string   `hcl:"using,optional"`
	Where        string   `hcl:"where,optional"`
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
