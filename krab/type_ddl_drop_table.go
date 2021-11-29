package krab

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
)

// DDLDropTable contains DSL for dropping tables.
type DDLDropTable struct {
	Name string `hcl:"name,label"`

	DefRange hcl.Range
}

var DDLDropTableSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "drop_table",
			LabelNames: []string{"name"},
		},
	},
}

// ToSQL converts migration definition to SQL.
func (d *DDLDropTable) ToSQL(w io.StringWriter) {
	w.WriteString("DROP TABLE ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
}
