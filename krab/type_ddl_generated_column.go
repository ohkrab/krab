package krab

import "io"

// DDLGeneratedColumn DSL.
type DDLGeneratedColumn struct {
	As string `hcl:"as"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLGeneratedColumn) ToSQL(w io.StringWriter) {
	w.WriteString("GENERATED ALWAYS AS (")
	w.WriteString(d.As)
	w.WriteString(") STORED")
}
