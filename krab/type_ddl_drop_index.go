package krab

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
)

// DDLDropIndex contains DSL for dropping indicies.
type DDLDropIndex struct {
	Name string `hcl:"name,label"`

	Cascade      bool `hcl:"cascade,optional"`
	Concurrently bool `hcl:"concurrently,optional"`

	DefRange hcl.Range
}

var DDLDropIndexSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "drop_index",
			LabelNames: []string{"name"},
		},
	},
}

// ToSQL converts migration definition to SQL.
func (d *DDLDropIndex) ToSQL(w io.StringWriter) {
	w.WriteString("DROP INDEX")
	if d.Concurrently {
		w.WriteString(" CONCURRENTLY")
	}
	w.WriteString(" ")
	w.WriteString(krabdb.QuoteIdentWithDots(d.Name))
	if d.Cascade {
		w.WriteString(" CASCADE")
	}
}
