package krab

import (
	"github.com/ohkrab/krab/krabhcl"
)

// TestExample represents test runner configuration.
//
type TestExample struct {
	TestSuiteRefName string           `hcl:"test_suite,label"`
	Name             string           `hcl:"name,label"`
	Its              []*TestExampleIt `hcl:"it,block"`
}

func (t *TestExample) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "test", Labels: []string{t.TestSuiteRefName, t.Name}}
}

func (t *TestExample) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(t.TestSuiteRefName),
	)
}

// TestExampleIt represents one use case for test example that contain assertions.
type TestExampleIt struct {
	Comment     string        `hcl:"comment,label"`
	Do          *Do           `hcl:"do,block"`
	RowAsserts  []*AssertRow  `hcl:"row,block"`
	RowsAsserts []*AssertRows `hcl:"rows,block"`
}

// AssertRows
type AssertRows struct {
	Expectations []*Expect `hcl:"expect,block"`
}

// AssertRow
type AssertRow struct {
	Scope        string    `hcl:"scope,label"`
	Expectations []*Expect `hcl:"expect,block"`
}

// Expect
type Expect struct {
	Subject string  `hcl:"subject,label"`
	Equal   *string `hcl:"equal,optional"`
}
