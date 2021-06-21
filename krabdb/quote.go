package krabdb

import (
	"github.com/jackc/pgx/v4"
)

func QuoteIdent(s string) string {
	return pgx.Identifier{s}.Sanitize()
}
