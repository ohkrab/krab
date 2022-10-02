package krab

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLDropIndex contains DSL for dropping indicies.
type DDLDropIndex struct {
	krabhcl.Source

	Name string `hcl:"name,label"`

	Cascade      bool `hcl:"cascade,optional"`
	Concurrently bool `hcl:"concurrently,optional"`
}

var DDLDropIndexSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "drop_index",
			LabelNames: []string{"name"},
		},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLDropIndex) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	panic("Not implemented drop index")
	d.Source.Extract(block)

	return nil
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
