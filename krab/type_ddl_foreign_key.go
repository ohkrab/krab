package krab

import (
	"io"
	"strings"

	"github.com/ohkrab/krab/krabdb"
)

// DDLForeignKey constraint DSL for table DDL.
type DDLForeignKey struct {
	Columns    []string      `hcl:"columns"`
	References DDLReferences `hcl:"references,block"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLForeignKey) ToSQL(w io.StringWriter) {
	w.WriteString("FOREIGN KEY (")
	cols := krabdb.QuoteIdentStrings(d.Columns)
	w.WriteString(strings.Join(cols, ","))
	w.WriteString(") ")
	d.References.ToSQL(w)
}
