package krab

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
)

// DDLCreateIndex contains DSL for creating indicies.
type DDLCreateIndex struct {
	krabhcl.Source

	Table string `hcl:"table,label"`
	Name  string `hcl:"name,label"`

	Unique       bool     `hcl:"unique,optional"`
	Concurrently bool     `hcl:"concurrently,optional"`
	Columns      []string `hcl:"columns"`
	Include      []string `hcl:"include,optional"`
	Using        string   `hcl:"using,optional"`
	Where        string   `hcl:"where,optional"`
}

var DDLCreateIndexSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "create_index",
			LabelNames: []string{"table", "name"},
		},
	},
}

// DecodeHCL parses HCL into struct.
func (d *DDLCreateIndex) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	panic("Not implemented create index")
	d.Source.Extract(block)

	d.Columns = []string{}
	d.Include = []string{}

	return nil
}

// ToSQL converts migration definition to SQL.
func (d *DDLCreateIndex) ToSQL(w io.StringWriter) {
	w.WriteString("CREATE")
	if d.Unique {
		w.WriteString(" UNIQUE")
	}
	w.WriteString(" INDEX")
	if d.Concurrently {
		w.WriteString(" CONCURRENTLY")
	}
	w.WriteString(" ")
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" ON ")
	w.WriteString(krabdb.QuoteIdent(d.Table))
	if d.Using != "" {
		w.WriteString(" USING ")
		w.WriteString(d.Using)
	}
	w.WriteString(" (")
	for i, col := range d.Columns {
		w.WriteString(krabdb.QuoteIdent(col))
		if i < len(d.Columns)-1 {
			w.WriteString(",")
		}
	}
	w.WriteString(")")

	if len(d.Include) > 0 {
		w.WriteString(" INCLUDE (")
		for i, col := range d.Include {
			w.WriteString(krabdb.QuoteIdent(col))
			if i < len(d.Columns)-1 {
				w.WriteString(",")
			}
		}
		w.WriteString(")")
	}

	if d.Where != "" {
		w.WriteString(" WHERE (")
		w.WriteString(d.Where)
		w.WriteString(")")
	}
}
