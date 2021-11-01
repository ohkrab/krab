package krabdb

import (
	"strings"

	"github.com/jackc/pgx/v4"
)

// QuoteIdent escapes identifiers in PG.
//
//     public -> "public"
//
func QuoteIdent(s string) string {
	return pgx.Identifier{s}.Sanitize()
}

// QuoteIdentWithDots escapes identifiers in PG.
//
//     public.test -> "public"."test"
//
func QuoteIdentWithDots(s string) string {
	names := strings.Split(s, ".")
	for i, name := range names {
		names[i] = QuoteIdent(name)
	}
	return strings.Join(names, ".")
}

// QuoteIdentStrings escapes identifiers in PG.
func QuoteIdentStrings(in []string) []string {
	out := make([]string, len(in))
	for i, name := range in {
		out[i] = QuoteIdent(name)
	}
	return out
}
