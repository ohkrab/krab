package krab

import (
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// ToSQL converts DSL struct to SQL.
type ToSQL interface {
	ToSQL(w io.StringWriter)
}

// SQLStatement represents raw SQL statement.
type SQLStatement string

// SQLStatements represents list of raw SQL statements.
type SQLStatements []SQLStatement

// SQLStatementsSorter sorts SQLStatement by the order how they are defined in a file.
type SQLStatementsSorter struct {
	Statements SQLStatements
	Bytes      []int
}

// Append adds new SQL statement to the list from object that satisfies ToSQL interface.
func (s *SQLStatements) Append(sql ToSQL) {
	sb := &strings.Builder{}
	sql.ToSQL(sb)
	*s = append(*s, SQLStatement(sb.String()))
}

// Insert ToSQL at given range.
func (s *SQLStatementsSorter) Insert(r hcl.Range, sql ToSQL) {
	s.Statements.Append(sql)
	s.Bytes = append(s.Bytes, r.Start.Byte)
}

// Sort sorts statements by byte range.
func (s *SQLStatementsSorter) Sort() SQLStatements {
	ret := make(SQLStatements, len(s.Statements))
	copy(ret, s.Statements) // TODO: replace with actual sort
	return ret
}
