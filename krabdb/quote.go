package krabdb

import "github.com/lib/pq"

func QuoteIdent(s string) string {
	return pq.QuoteIdentifier(s)
}
