package krab

import "io"

// DDLIdentity DSL.
type DDLIdentity struct {
	// Generated string `hcl:"generated,optional"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLIdentity) ToSQL(w io.StringWriter) {
	w.WriteString("GENERATED ALWAYS AS IDENTITY")
}
