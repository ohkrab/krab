package krabtpl

import (
	"text/template"

	"github.com/ohkrab/krab/krabdb"
)

var Functions template.FuncMap = template.FuncMap{
	"quote": krabdb.QuoteIdent,
}
