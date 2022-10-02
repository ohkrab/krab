package krab

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLDropTable contains DSL for dropping tables.
type DDLDropTable struct {
	krabhcl.Source

	Name string `hcl:"name,label"`
}

var DDLDropTableSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "drop_table",
			LabelNames: []string{"name"},
		},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLDropTable) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	panic("Not implemented Drop Table")
	d.Source.Extract(block)

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLDropTable) ToSQL(w io.StringWriter) {
	w.WriteString("DROP TABLE ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
}
