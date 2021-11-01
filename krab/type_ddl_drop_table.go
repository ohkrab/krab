package krab

import (
	"io"

	"github.com/ohkrab/krab/krabdb"
)

// DDLDropTable contains DSL for dropping tables.
type DDLDropTable struct {
	Name string `hcl:"name,label"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLDropTable) ToSQL(w io.StringWriter) {
	w.WriteString("DROP TABLE ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
}
