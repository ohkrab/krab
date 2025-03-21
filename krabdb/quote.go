package krabdb

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// QuoteIdent escapes identifiers in PG.
//
//	public -> "public"
func QuoteIdent(s string) string {
	return pgx.Identifier{s}.Sanitize()
}

// QuoteIdentWithDots escapes identifiers in PG.
//
//	public.test -> "public"."test"
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

// Quote escapes values in PG.
func Quote(o any) string {
	switch o := o.(type) {
	case nil:
		return "null"
	case int64:
		return strconv.FormatInt(o, 10)
	case uint64:
		return strconv.FormatUint(o, 10)
	case float64:
		return strconv.FormatFloat(o, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(o)
	case []byte:
		return `'\x` + hex.EncodeToString(o) + "'"
	case string:
		return "'" + strings.ReplaceAll(o, "'", "''") + "'"
	case time.Time:
		return o.Truncate(time.Microsecond).Format("'2006-01-02 15:04:05.999999999Z07:00:00'")
	default:
		panic(fmt.Sprintf("Quote not implemented for type %T", o))
	}
}
