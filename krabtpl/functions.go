package krabtpl

import (
	"text/template"

	"github.com/jaswdr/faker"
	"github.com/ohkrab/krab/krabdb"
)

func Functions() template.FuncMap {
	fake := faker.New()

	return template.FuncMap{
		"quote_ident": krabdb.QuoteIdent,
		"quote":       krabdb.Quote,
		"fake":        Fake(fake),
	}
}
