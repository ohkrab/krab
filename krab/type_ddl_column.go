package krab

import (
	"fmt"
	"io"
	"strings"

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
			w.WriteString(krabdb.Quote(defaultExpr.AsBool()))

		case cty.Number:
			switch strings.ToLower(d.Type) {
			case "smallint", "integer", "int", "bigint", "smallserial", "serial", "bigserial":
				w.WriteString(krabdb.Quote(defaultExpr.AsInt64()))
			case "real", "double precision":
				w.WriteString(krabdb.Quote(defaultExpr.AsFloat64()))
			default:
				//TODO: implement big numbers (numeric, decimal)
				panic(fmt.Sprintf(
					"%sCannot map default type of %s to SQL, if you see this error please report the issue with example so I can fix this%s",
					ctc.BackgroundRed|ctc.ForegroundYellow,
					d.Type,
					ctc.Reset,
				))
			}

		case cty.String:
			w.WriteString(krabdb.Quote(defaultExpr.AsString()))

		default:
			switch {
			case defaultExpr.Type().IsObjectType():
				w.WriteString("'{}'")

			case defaultExpr.Type().IsTupleType():
				w.WriteString("'[]'")

			default:
				panic(fmt.Sprintf(
					"%sCannot map default type %s to SQL, if you see this error please report the issue with example so I can fix this%s",
					ctc.BackgroundRed|ctc.ForegroundYellow,
					d.Type,
					ctc.Reset,
				))
			}
		}
	}
}
