package krab

import (
	"io"
	"sort"
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

// Append adds new SQL statement to the list from object that satisfies ToSQL interface.
func (s *SQLStatements) Append(sql ToSQL) {
	sb := &strings.Builder{}
	sql.ToSQL(sb)
	*s = append(*s, SQLStatement(sb.String()))
}

// SQLStatementsSorter sorts SQLStatement by the order how they are defined in a file.
type SQLStatementsSorter struct {
	Statements SQLStatements
	Bytes      []int
}

// Len is the number of elements in the collection.
func (s *SQLStatementsSorter) Len() int {
	return len(s.Statements)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (s *SQLStatementsSorter) Less(i int, j int) bool {
	return s.Bytes[i] < s.Bytes[j]
}

// Swap swaps the elements with indexes i and j.
func (s *SQLStatementsSorter) Swap(i int, j int) {
	s.Bytes[i], s.Bytes[j] = s.Bytes[j], s.Bytes[i]
	s.Statements[i], s.Statements[j] = s.Statements[j], s.Statements[i]
}

// Insert ToSQL at given range.
func (s *SQLStatementsSorter) Insert(r hcl.Range, sql ToSQL) {
	s.Statements.Append(sql)
	s.Bytes = append(s.Bytes, r.Start.Byte)
}

// Sort sorts statements by byte range.
func (s *SQLStatementsSorter) Sort() SQLStatements {
	sort.Sort(s)
	return s.Statements
}
