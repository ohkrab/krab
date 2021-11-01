package krab

import (
	"io"

	"github.com/ohkrab/krab/krabdb"
)

// DDLCheck constraint DSL for table DDL.
type DDLCheck struct {
	Name       string `hcl:"name,label"`
	Expression string `hcl:"expression"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLCheck) ToSQL(w io.StringWriter) {
	w.WriteString("CONSTRAINT ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" CHECK (")
	w.WriteString(d.Expression)
	w.WriteString(")")
}
