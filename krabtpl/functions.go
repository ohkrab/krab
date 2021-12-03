package krabtpl

import (
	"text/template"

	"github.com/ohkrab/krab/krabdb"
)

var Functions template.FuncMap = template.FuncMap{
	"quote_ident": krabdb.QuoteIdent,
	"quote":       krabdb.Quote,
}
