package krab

import (
	"io"
	"strings"

	"github.com/ohkrab/krab/krabdb"
)

// DDLUnique constraint DSL for table DDL.
type DDLUnique struct {
	Columns []string `hcl:"columns"`
	Include []string `hcl:"include,optional"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLUnique) ToSQL(w io.StringWriter) {
	w.WriteString("UNIQUE (")
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
