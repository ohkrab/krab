package krabdb

import (
	"github.com/jackc/pgx/v4"
)

// QuoteIdent escapes identifiers in PG.
func QuoteIdent(s string) string {
	return pgx.Identifier{s}.Sanitize()
}
