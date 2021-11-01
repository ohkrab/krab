package krab

import (
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/wzshiming/ctc"
	"github.com/zclconf/go-cty/cty"
)

// DDLColumn DSL for table DDL.
type DDLColumn struct {
	Name      string              `hcl:"name,label"`
	Type      string              `hcl:"type,label"`
	Null      *bool               `hcl:"null,optional"`
	Identity  *DDLIdentity        `hcl:"identity,block"`
	Default   hcl.Expression      `hcl:"default,optional"`
	Generated *DDLGeneratedColumn `hcl:"generated,block"`
}

// ToSQL converts migration definition to SQL.
func (d *DDLColumn) ToSQL(w io.StringWriter) {
	w.WriteString(krabdb.QuoteIdent(d.Name))
	w.WriteString(" ")
	w.WriteString(d.Type)

	if d.Null != nil {
		w.WriteString(" ")
		if *d.Null {
			w.WriteString("NULL")
		} else {
			w.WriteString("NOT NULL")
		}
	}

	if d.Identity != nil {
		w.WriteString(" ")
		d.Identity.ToSQL(w)
	}

	if d.Generated != nil {
		w.WriteString(" ")
		d.Generated.ToSQL(w)
	}

	defaultExpr := krabhcl.Expression{Expr: d.Default}
	if defaultExpr.Ok() {
		w.WriteString(" DEFAULT ")

		switch defaultExpr.Type() {
		case cty.Bool:
			w.WriteString(strconv.FormatBool(defaultExpr.AsBool()))

		default:
			panic(fmt.Sprint(
				ctc.BackgroundRed|ctc.ForegroundYellow,
				"Cannot map default type to SQL, if you see this error please report the issue with example",
				ctc.Reset,
			))
		}
	}
}
