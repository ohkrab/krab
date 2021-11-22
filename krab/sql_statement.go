package krab

import (
	"io"
	"strings"
)

// ToSQL converts DSL struct to SQL.
type ToSQL interface {
	ToSQL(w io.StringWriter)
}

// SQLStatement represents raw SQL statement.
type SQLStatement string

// SQLStatements represents list of raw SQL statements.
type SQLStatements []SQLStatement

// Append adds new SQL statement to the list from object that satisfies ToSQL interface.
func (s *SQLStatements) Append(sql ToSQL) {
	sb := &strings.Builder{}
	sql.ToSQL(sb)
	*s = append(*s, SQLStatement(sb.String()))
}
