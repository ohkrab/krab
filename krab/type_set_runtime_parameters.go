package krab

import (
	"io"
)

// SetRuntimeParameters
// https://www.postgresql.org/docs/current/sql-set.html
type SetRuntimeParameters struct {
	SearchPath *string `hcl:"search_path"`
}

// ToSQL converts set parameters to SQL.
func (d *SetRuntimeParameters) ToSQL(w io.StringWriter) {
	if d.SearchPath != nil {
		w.WriteString("SET search_path TO ")
		w.WriteString(*d.SearchPath)
	}
}
